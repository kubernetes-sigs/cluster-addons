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

package coredns

import (
	"context"
	"strings"

	api "sigs.k8s.io/addon-operators/coredns/pkg/apis/addons/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/status"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"
)

// Add creates a new CoreDNS Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileCoreDNS {
	labels := map[string]string{
		"k8s-app": "coredns",
	}

	r := &ReconcileCoreDNS{}

	replacePlaceholders := func(ctx context.Context, object declarative.DeclarativeObject, s string) (string, error) {
		// TODO: Should we default and if so where?
		dnsDomain := "" // o.Spec.DNSDomain
		if dnsDomain == "" {
			dnsDomain = "cluster.local"
		}

		dnsServerIP := "" // o.Spec.DNSServerIP
		if dnsServerIP == "" {
			dnsServerIP = "10.0.0.10"
		}

		s = strings.Replace(s, "__PILLAR__DNS__DOMAIN__", dnsDomain, -1)
		s = strings.Replace(s, "__PILLAR__DNS__SERVER__", dnsServerIP, -1)

		return s, nil
	}

	r.Reconciler.Init(mgr, &api.CoreDNS{},
		declarative.WithRawManifestOperation(replacePlaceholders),
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(declarative.SourceLabel(mgr.GetScheme())),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		declarative.WithApplyPrune(),
	)

	return r
}

func add(mgr manager.Manager, r *ReconcileCoreDNS) error {
	// Create a new controller
	c, err := controller.New("coredns-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to CoreDNS
	err = c.Watch(&source.Kind{Type: &api.CoreDNS{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to deployed objects
	_, err = declarative.WatchAll(mgr.GetConfig(), c, r, declarative.SourceLabel(mgr.GetScheme()))
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCoreDNS{}

// +kubebuilder:rbac:groups=addons.k8s.io,resources=coredns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.k8s.io,resources=coredns/status,verbs=get;update;patch

// ReconcileCoreDNS reconciles a CoreDNS object
type ReconcileCoreDNS struct {
	declarative.Reconciler
}
