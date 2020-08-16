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

// The cluster-addons-bootstrap ensures that a given set of addon Objects is installed in the cluster.
package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-addons/bootstrap/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	klog.InitFlags(nil)
	flag.Set("alsologtostderr", "false") // false is default, but this is informative

	flag.Parse()
	// make sure we flush before exiting
	defer klog.Flush()

	am, err := app.AddonManager(os.Getenv)

	if err != nil {
		klog.Exit(err)
	}

	if err := am.Run(); err != nil {
		klog.Exit(err)
	}
}
