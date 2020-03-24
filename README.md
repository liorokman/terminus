# Terminus


Terminus is a DaemonSet which allows setting pod-level resource limits on Kubernetes pods.

This project was created as part of the validation work for [KEP 1592](https://github.com/kubernetes/enhancements/pull/1592).

# Installation

## Configuration

1. Edit `config/manager/container_runtime_endpoint.yaml` and set the remote-cri socket path correctly for your cluster. If your cluster is using Docker greater or equal to version 18, then Docker is actually running containerd for managing its containers. You can find the correct path by running `ps axf| grep shim`, or checking `/var/run/docker/containerd`
1. Build the image by running `make docker-image IMG=terminus:1`.  Use an `IMG` value that you can pull from your cluster.
1. Push the newly build image (`terminus:1`) to somewhere that your cluster can pull it from. If running on `kind`, then use `kind load docker-image terminus:1`.
1. Deploy Terminus to your cluster by running `make deploy IMG=terminus:1`

## Annotate pods

Add an annotation called `podlimits.terminus.sap.com/pod-limits` to pods which should be managed by the node-local controller.

The value of the annotation is a JSON document containing a standard Kubernetes `ResourceList` object under the `limits` key. For example:

```json
{ 
  "limits": {
     "memory":"64M",
     "cpu": "100m",
  } 
}
```

# Effect

When deployed, terminus starts a Pod watcher on all of the nodes where it is running. Whenever a Pod reaches `Running` state, if it is annotated with a `podlimits.terminus.sap.com/pod-limits` annotation the pod-level cgroup is adjusted as required. For CPU, the pod level quota is set to the provided limit. For memory, the pod level `mem.limit_in_bytes` is set to the provided limit and all cgroups for all containers that provide a request value are updated so that the soft limit (`memory.soft_limit_in_bytes`) is set to the request value for the container.

No validation is made to verify that the pod level limits are larger or equal to container level requests, and no limits are removed form speific containers. 
