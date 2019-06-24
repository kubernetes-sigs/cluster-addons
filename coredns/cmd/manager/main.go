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

package main

import (
	"flag"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	"k8s.io/klog/klogr"
	"sigs.k8s.io/addon-operators/coredns/pkg/apis"
	"sigs.k8s.io/addon-operators/coredns/pkg/controller"
	"sigs.k8s.io/addon-operators/coredns/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
)

func main() {
	klog.InitFlags(nil)
	addon.Init()

	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()
	logf.SetLogger(klogr.New())

	// Get a config to talk to the apiserver
	klog.Info("setting up client for manager")
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	klog.Info("setting up manager")
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: metricsAddr})
	if err != nil {
		klog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	klog.Info("Registering Components.")

	// Setup Scheme for all resources
	klog.Info("setting up scheme")
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	// Setup all Controllers
	klog.Info("Setting up controller")
	if err := controller.AddToManager(mgr); err != nil {
		klog.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	klog.Info("setting up webhooks")
	if err := webhook.AddToManager(mgr); err != nil {
		klog.Error(err, "unable to register webhooks to the manager")
		os.Exit(1)
	}

	// Start the Cmd
	klog.Info("Starting the Cmd.")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Error(err, "unable to run the manager")
		os.Exit(1)
	}
}
