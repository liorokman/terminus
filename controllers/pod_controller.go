/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"

	"github.com/go-logr/logr"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	terminusruntime "github.com/liorokman/terminus/runtime"
)

const (
	limitAnnotation string = "podlimits.terminus.sap.com/pod-limits"
)

// PodReconciler reconciles a Pod object
type PodReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type LimitAnnotation struct {
	Limits corev1.ResourceList `json:"limits,omitempty"`
}

// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods/status,verbs=get

func (r *PodReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("pod", req.NamespacedName)

	var thepod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &thepod); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "unable to fetch the pod")
		}
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var podLimits LimitAnnotation
	if thepod.Status.HostIP != viper.GetString("host.ip") {
		return ctrl.Result{}, nil
	} else if thepod.Status.Phase != corev1.PodRunning {
		return ctrl.Result{}, nil
	} else if limit, ok := thepod.GetAnnotations()[limitAnnotation]; !ok {
		return ctrl.Result{}, nil
	} else if err := json.Unmarshal([]byte(limit), &podLimits); err != nil {
		log.Error(err, "Failed to unmarshal limit annotation", limitAnnotation, limit)
		return ctrl.Result{}, err
	}
	log.Info("Pod received", "pod", thepod)

	// Set the soft-limit correctly on a per-container basis
	for _, c := range thepod.Spec.Containers {
		if req, ok := c.Resources.Requests[corev1.ResourceMemory]; ok {
			log.Info("Setting reservation limit", "name", c.Name, "request-value", req)
			for _, cs := range thepod.Status.ContainerStatuses {
				if c.Name == cs.Name {
					reqInBytes := req.Value()
					if cgroup, err := terminusruntime.LoadCgroup(ctx, cs.ContainerID, false); err == nil {
						err := cgroup.Update(&specs.LinuxResources{
							Memory: &specs.LinuxMemory{
								Reservation: &reqInBytes,
							},
						})
						if err != nil {
							log.Error(err, "Failed to set soft-limit for a container", "name", c.Name)
						}
					} else {
						log.Error(err, "Failed to load the container cgroup", "name", c.Name, "containerID", cs.ContainerID)
					}
				}
			}
		}
	}

	// Now fix the top-level cgroups
	cgroup, err := terminusruntime.LoadCgroup(ctx, thepod.Status.ContainerStatuses[0].ContainerID, true)
	if err != nil {
		log.Error(err, "Failed to find the pod cgroup")
		return ctrl.Result{}, err
	}
	log.Info("Cgroup values", "cgroup", cgroup)
	resourceRequest := specs.LinuxResources{}
	if quant, ok := podLimits.Limits[corev1.ResourceMemory]; ok {
		lim := quant.Value()
		resourceRequest.Memory = &specs.LinuxMemory{
			Limit: &lim,
		}
	}
	if quant, ok := podLimits.Limits[corev1.ResourceCPU]; ok {
		quota := terminusruntime.MilliCPUToQuota(quant.Value())
		period := uint64(terminusruntime.QuotaPeriod)
		resourceRequest.CPU = &specs.LinuxCPU{
			Quota:  &quota,
			Period: &period,
		}
	}
	log.Info("Setting pod level cgroups", "ResourceRequest", resourceRequest)
	err = cgroup.Update(&resourceRequest)
	if err != nil {
		log.Error(err, "Failed to update the pod limit")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PodReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
