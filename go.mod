module github.com/liorokman/terminus

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5 // indirect
	github.com/Microsoft/hcsshim v0.8.6 // indirect
	github.com/containerd/cgroups v0.0.0-20200308110149-6c3dec43a103
	github.com/containerd/containerd v1.3.2
	github.com/containerd/continuity v0.0.0-20200228182428-0f16d7a0959c // indirect
	github.com/containerd/fifo v0.0.0-20191213151349-ff969a566b00 // indirect
	github.com/containerd/ttrpc v0.0.0-20200121165050-0be804eadb15 // indirect
	github.com/containerd/typeurl v0.0.0-20200205145503-b45ef1f1f737 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/go-events v0.0.0-20190806004212-e31b211e4f1c // indirect
	github.com/go-logr/logr v0.1.0
	github.com/gogo/googleapis v1.3.2 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v1.0.0-rc10 // indirect
	github.com/opencontainers/runtime-spec v0.1.2-0.20190507144316-5b71a03e2700
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/spf13/viper v1.6.2
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	google.golang.org/genproto v0.0.0-20200117163144-32f20d992d24 // indirect
	google.golang.org/grpc v1.26.0 // indirect
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	sigs.k8s.io/controller-runtime v0.4.0
)

replace github.com/docker/distribution v2.7.1+incompatible => github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible
