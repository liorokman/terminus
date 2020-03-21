package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/containerd/cgroups"
	"github.com/containerd/containerd"
)

func main() {
	var path string
	var cid string
	var namespace string

	flag.StringVar(&path, "address", "/var/run/docker/containerd/containerd.sock", "address to containerd socket")
	flag.StringVar(&cid, "container", "", "container id to query")
	flag.StringVar(&namespace, "namespace", "moby", "container id to query")
	flag.Parse()

	client, err := containerd.New(path, containerd.WithDefaultNamespace(namespace))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(-1)
	}
	containers, err := client.Containers(context.TODO())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(-1)
	}
	for _, c := range containers {
		if c.ID() != cid {
			continue
		}
		info, err := c.Spec(context.TODO())
		if err != nil {
			fmt.Printf("Failed: %s\n", err.Error())
			continue
		}
		fmt.Printf("%#v\n", info.Linux.CgroupsPath)
		control, err := cgroups.Load(cgroups.V1, cgroups.StaticPath(info.Linux.CgroupsPath))
		if err != nil {
			fmt.Printf("Failed: %s\n", err.Error())
			continue
		}
		metrics, err := control.Stat(func(e error) error {
			fmt.Printf("Error generating stats: %s\n", e.Error())
			return nil
		})
		if err != nil {
			fmt.Printf("Failed: %s\n", err.Error())
			continue
		}
		fmt.Printf("memory: %#v\n", metrics.Memory)
	}

}
