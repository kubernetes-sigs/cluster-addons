/*
Copyright 2020 The Kubernetes Authors.

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
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	api "sigs.k8s.io/cluster-addons/coredns/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/status"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"

	"github.com/pkg/errors"
)

var _ reconcile.Reconciler = &CoreDNSReconciler{}

// CoreDNSReconciler reconciles a CoreDNS object
type CoreDNSReconciler struct {
	Client      client.Client
	Scheme      *runtime.Scheme
	watchLabels declarative.LabelMaker

	declarative.Reconciler
}

func (r *CoreDNSReconciler) setupReconciler(mgr ctrl.Manager) error {
	labels := map[string]string{
		"k8s-app": "kube-dns",
	}

	replacePlaceholders := func(ctx context.Context, object declarative.DeclarativeObject, s string) (string, error) {
		o := object.(*api.CoreDNS)

		if o.Spec.DNSDomain == "" {
			domain := getDNSDomain()
			o.Spec.DNSDomain = domain
		}
		if o.Spec.DNSIP == "" {
			ip, err := findDNSClusterIP(ctx, mgr.GetClient())
			if err != nil {
				return "", errors.Errorf("unable to find CoreDNS IP: %v", err)
			}
			o.Spec.DNSIP = ip
		}
		s = strings.Replace(s, "{{ .DNSDomain }}", o.Spec.DNSDomain, -1)
		s = strings.Replace(s, "{{ .DNSIP }}", o.Spec.DNSIP, -1)
		return s, nil
	}

	replaceCorefilePlaceholder := func(ctx context.Context, object declarative.DeclarativeObject, s string) (string, error) {
		var err error
		var corefile string

		o := object.(*api.CoreDNS)
		if o.Spec.Corefile == "" {
			corefile, err = getCorefile(ctx, mgr.GetClient())
			if err != nil {
				return "", errors.Errorf("unable to find CoreDNS Corefile: %v", err)
			}
			o.Spec.Corefile = corefile
		}

		// Check for Corefile Migration
		corefile, err = corefileMigration(ctx, mgr.GetClient(), o.Spec.Version, o.Spec.Corefile)
		if err != nil {
			return "", err
		}

		// Usually returns an empty Corefile if the Corefile is default.
		if corefile == "" {
			corefilePath := fmt.Sprintf("channels/packages/coredns/%s/Corefile", o.Spec.Version)
			b, err := ioutil.ReadFile(corefilePath)
			if err != nil {
				return "", err
			}
			corefile = string(b)
		}

		corefile = strings.Replace(corefile, "{{ .DNSDomain }}", o.Spec.DNSDomain, -1)
		o.Spec.Corefile = corefile

		s = strings.Replace(s, "{{ .Corefile }}", prepCorefileFormat(o.Spec.Corefile, 8), -1)

		return s, nil
	}

	r.watchLabels = declarative.SourceLabel(mgr.GetScheme())

	return r.Reconciler.Init(mgr, &api.CoreDNS{},
		declarative.WithRawManifestOperation(replaceCorefilePlaceholder),
		declarative.WithRawManifestOperation(replacePlaceholders),
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(r.watchLabels),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		declarative.WithObjectTransform(addon.ApplyPatches),
		declarative.WithApplyPrune(),
		declarative.WithApplyKustomize(),
	)
}

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=coredns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=coredns/status,verbs=get;update;patch

func (r *CoreDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addon.Init()

	if err := r.setupReconciler(mgr); err != nil {
		return err
	}

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
	_, err = declarative.WatchAll(mgr.GetConfig(), c, r, r.watchLabels)
	if err != nil {
		return err
	}

	return nil
}

// for WithApplyPrune
// +kubebuilder:rbac:groups=*,resources=*,verbs=list

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=coredns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=coredns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps;extensions,resources=deployments,verbs=get;list;watch;create;update;delete;patch
// +kubebuilder:rbac:groups="",resources=configmaps;serviceaccounts;services,verbs=get;list;watch;create;update;delete;patch
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings;clusterroles;clusterrolebindings,verbs=get;list;watch;create;update;delete;patch
// To grant permissions to CoreDNS, we need those permissions:
// +kubebuilder:rbac:groups="",resources=endpoints;namespaces;nodes;pods,verbs=get;list;watch
