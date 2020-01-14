/*
Copyright 2020 The Kubernetes authors.

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
	"fmt"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage/names"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	// +kubebuilder:scaffold:imports

	discoveryv1alpha1 "sigs.k8s.io/cluster-addons/discovery/api/v1alpha1"
	"sigs.k8s.io/cluster-addons/discovery/controllers/decorators"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const (
	timeout  = time.Second * 10
	interval = time.Millisecond * 100
)

var (
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	scheme    = runtime.NewScheme()

	gracePeriod int64 = 0
	propagation       = metav1.DeletePropagationForeground
	deleteOpts        = &client.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
		PropagationPolicy:  &propagation,
	}
	genName = names.SimpleNameGenerator.GenerateName
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = discoveryv1alpha1.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())
	err = k8sscheme.AddToScheme(scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	By("Setting up a controller manager")
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())

	reconciler, err := NewAddonReconciler(
		mgr.GetClient(),
		ctrl.Log.WithName("controllers").WithName("Addon"),
		mgr.GetScheme(),
	)
	Expect(err).ToNot(HaveOccurred())

	By("Adding the addon controller to the manager")
	Expect(reconciler.SetupWithManager(mgr)).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()

		By("Starting managed controllers")
		err = mgr.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = mgr.GetClient()
	Expect(k8sClient).ToNot(BeNil())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func newAddon(name string) *decorators.Addon {
	return &decorators.Addon{
		Addon: &discoveryv1alpha1.Addon{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	}
}

func toRefs(scheme *runtime.Scheme, objs ...runtime.Object) (refs []discoveryv1alpha1.RichReference) {
	for _, obj := range objs {
		ref, err := reference.GetReference(scheme, obj)
		if err != nil {
			panic(fmt.Errorf("error creating resource reference: %v", err))
		}

		// Clear unnecessary fields
		ref.UID = ""
		ref.ResourceVersion = ""

		refs = append(refs, discoveryv1alpha1.RichReference{
			ObjectReference: ref,
		})
	}

	return
}
