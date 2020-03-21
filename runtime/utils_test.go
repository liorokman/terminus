package runtime

import (
	"fmt"
	"testing"
)

func TestCgroupsPathToStaticPath(t *testing.T) {

	res := cgroupsPathToStaticPath("kubepods-burstable-pod249c6d91_68eb_11e9_b001_ba9d0aa05f94.slice:docker:94d70466d4b3139d02fb9a3e1af1624dcad732f7795fc2a7509407abfa3df434", true)
	fmt.Println(res)
	res = cgroupsPathToStaticPath("kubepods-burstable-pod249c6d91_68eb_11e9_b001_ba9d0aa05f94.slice:docker:94d70466d4b3139d02fb9a3e1af1624dcad732f7795fc2a7509407abfa3df434", false)
	fmt.Println(res)
	res = cgroupsPathToStaticPath("/kubepods/besteffort/pod1423c184-fc22-46cb-afb8-aa5bba016b71/5eb7dbec4ca0eafa656252868970250faa51fe67f15494fb99f2f012daccfa71", true)
	fmt.Println(res)
	res = cgroupsPathToStaticPath("/kubepods/besteffort/pod1423c184-fc22-46cb-afb8-aa5bba016b71/5eb7dbec4ca0eafa656252868970250faa51fe67f15494fb99f2f012daccfa71", false)
	fmt.Println(res)
}
