package runtime

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/containerd/cgroups"
	"github.com/containerd/containerd"
	"github.com/spf13/viper"
)

var ContainerNotFoundError error = errors.New("Container wasn't found")

func LoadCgroup(ctx context.Context, containerID string, podLevel bool) (cgroups.Cgroup, error) {
	namespace := ""
	parts := strings.SplitN(containerID, "://", 2)
	switch parts[0] {
	case "containerd":
		namespace = "k8s.io"
	case "docker":
		namespace = "moby"
	default:
		return nil, errors.New("terminus: Unrecognized container interface")
	}
	client, err := containerd.New(viper.GetString("container.runtime.endpoint"), containerd.WithDefaultNamespace(namespace))
	if err != nil {
		return nil, err
	}
	containers, err := client.Containers(ctx, "id == "+parts[1])
	if err != nil {
		return nil, err
	}
	if len(containers) != 1 {
		return nil, ContainerNotFoundError
	}
	info, err := containers[0].Spec(ctx)
	if err != nil {
		return nil, err
	}
	return cgroups.Load(cgroups.V1, cgroups.StaticPath(cgroupsPathToStaticPath(info.Linux.CgroupsPath, podLevel)))
}

const systemdSuffix string = ".slice"

func cgroupsPathToStaticPath(path string, podLevel bool) string {
	if !strings.Contains(path, systemdSuffix) {
		if !podLevel {
			return path
		}
		return filepath.Dir(path)
	}
	parts := strings.Split(path, ":")
	if len(parts) < 1 {
		return path
	}
	segments := strings.Split(parts[0], "-")
	result := make([]string, len(segments))
	prevPrefix := ""
	for i := range segments {
		if i == 0 {
			result[i] = "/" + segments[i] + systemdSuffix
			prevPrefix = segments[i]
		} else {
			result[i] = prevPrefix + "-" + segments[i]
			if !strings.HasSuffix(result[i], systemdSuffix) {
				result[i] = result[i] + systemdSuffix
			}
			prevPrefix = prevPrefix + "-" + segments[i]
		}
	}
	if !podLevel {
		if len(parts) == 3 {
			result = append(result, parts[1]+"-"+parts[2]+".scope")
		} else {
			return path
		}
	}
	return strings.Join(result, "/")
}

const (
	// Taken from lmctfy https://github.com/google/lmctfy/blob/master/lmctfy/controllers/cpu_controller.cc
	MilliCPUToCPU int64 = 1000

	// 100000 is equivalent to 100ms
	QuotaPeriod    int64 = 100000
	MinQuotaPeriod int64 = 1000
)

// Copied with care from the Kubelet source code.
// MilliCPUToQuota converts milliCPU to CFS quota and period values.
func MilliCPUToQuota(milliCPU int64) (quota int64) {
	// CFS quota is measured in two values:
	//  - cfs_period_us=100ms (the amount of time to measure usage across given by period)
	//  - cfs_quota=20ms (the amount of cpu time allowed to be used across a period)
	// so in the above example, you are limited to 20% of a single CPU
	// for multi-cpu environments, you just scale equivalent amounts
	// see https://www.kernel.org/doc/Documentation/scheduler/sched-bwc.txt for details

	if milliCPU == 0 {
		return
	}

	// we then convert your milliCPU to a value normalized over a period
	quota = (milliCPU * QuotaPeriod) / MilliCPUToCPU

	// quota needs to be a minimum of 1ms.
	if quota < MinQuotaPeriod {
		quota = MinQuotaPeriod
	}
	return
}
