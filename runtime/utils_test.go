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

package runtime

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
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

func TestMilliCPUToQuota(t *testing.T) {

	quant := resource.MustParse("500m")
	res := MilliCPUToQuota(quant.MilliValue())
	fmt.Println(res)
}
