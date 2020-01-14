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
	"context"
	"sync"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	discoveryv1alpha1 "sigs.k8s.io/cluster-addons/discovery/api/v1alpha1"
	"sigs.k8s.io/cluster-addons/discovery/controllers/decorators"
	"sigs.k8s.io/cluster-addons/discovery/lib/controller-runtime/source"
)

// AddonReconciler reconciles a Addon object.
type AddonReconciler struct {
	client.Client

	log     logr.Logger
	mu      sync.RWMutex
	factory decorators.AddonFactory

	// addons contains the names of Addons the AddonReconciler has observed exist.
	addons map[types.NamespacedName]struct{}
	source *source.Dynamic
}

// +kubebuilder:rbac:groups=discovery.addons.x-k8s.io,resources=addons,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups=discovery.addons.x-k8s.io,resources=addons/status,verbs=update;patch
// +kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch

// SetupWithManager adds the addon reconciler to the given controller manager.
func (r *AddonReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Trigger addon events from the events of their compoenents.
	enqueueAddon := &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(r.mapComponentRequests),
	}

	// Add reconciler enqueued by dynamic Source watching all GVKs.
	return ctrl.NewControllerManagedBy(mgr).
		For(&discoveryv1alpha1.Addon{}).
		Watches(r.source, enqueueAddon).
		Complete(r)
}

// NewAddonReconciler constructs and returns an AddonReconciler.
// As a side effect, the given scheme has addon discovery types added to it
func NewAddonReconciler(cli client.Client, log logr.Logger, scheme *runtime.Scheme) (*AddonReconciler, error) {
	factory, err := decorators.NewSchemedAddonFactory(scheme)
	if err != nil {
		return nil, err
	}

	return &AddonReconciler{
		Client: cli,

		log:     log,
		factory: factory,
		addons:  map[types.NamespacedName]struct{}{},
		source:  &source.Dynamic{},
	}, nil
}

// Implement reconcile.Reconciler so the controller can reconcile objects
var _ reconcile.Reconciler = &AddonReconciler{}

// Reconcile transitions the state of an Operator resource based on the current state of the cluster.
func (r *AddonReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// Set up a convenient log object so we don't have to type request over and over again
	log := r.log.WithValues("request", req)
	log.V(1).Info("reconciling addon")

	// Fetch the Addon from the cache
	ctx := context.TODO()
	in := &discoveryv1alpha1.Addon{}
	if err := r.Get(ctx, req.NamespacedName, in); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Could not find Addon")
			r.unobserve(req.NamespacedName)
			// TODO(njhale): Recreate addon if we can find any components.
		} else {
			log.Error(err, "Error finding Addon")
		}

		return reconcile.Result{}, nil
	}
	r.observe(req.NamespacedName)

	// Wrap with convenience decorator
	addon, err := r.factory.NewAddon(in)
	if err != nil {
		log.Error(err, "Could not wrap Addon with convenience decorator")
		return reconcile.Result{}, nil
	}

	if err = r.updateComponents(ctx, addon); err != nil {
		log.Error(err, "Could not update components")
		return reconcile.Result{}, nil

	}

	if err := r.Update(ctx, addon.Addon); err != nil {
		log.Error(err, "Could not update Addon status")
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, req.NamespacedName, addon.Addon); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *AddonReconciler) updateComponents(ctx context.Context, addon *decorators.Addon) error {
	selector, err := addon.ComponentSelector()
	if err != nil {
		return err
	}

	components, err := r.listComponents(ctx, selector)
	if err != nil {
		return err
	}

	return addon.SetComponents(components...)
}

func (r *AddonReconciler) listComponents(ctx context.Context, selector labels.Selector) ([]runtime.Object, error) {
	informable, err := r.source.InformableGVKs()
	if err != nil {
		return nil, err
	}

	var componentLists []runtime.Object
	for _, gvk := range informable {
		gvk.Kind = gvk.Kind + "List"
		ul := &unstructured.UnstructuredList{}
		ul.SetGroupVersionKind(gvk)
		componentLists = append(componentLists, ul)
	}

	opt := client.MatchingLabelsSelector{Selector: selector}
	for _, list := range componentLists {
		if err := r.List(ctx, list, opt); err != nil {
			return nil, err
		}
	}

	return componentLists, nil
}

func (r *AddonReconciler) observed(name types.NamespacedName) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.addons[name]
	return ok
}

func (r *AddonReconciler) observe(name types.NamespacedName) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.addons[name] = struct{}{}
}

func (r *AddonReconciler) unobserve(name types.NamespacedName) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.addons, name)
}

func (r *AddonReconciler) mapComponentRequests(obj handler.MapObject) (requests []reconcile.Request) {
	if obj.Meta == nil {
		return
	}

	for _, name := range decorators.AddonNames(obj.Meta.GetLabels()) {
		// Only enqueue if we can find the addon in our cache
		if r.observed(name) {
			requests = append(requests, reconcile.Request{NamespacedName: name})
			continue
		}

		// Otherwise, best-effort generate a new addon
		// TODO(njhale): Implement verification that the addon-discovery admission webhook accepted this label (JWT or maybe sign a set of fields?)
		addon := &discoveryv1alpha1.Addon{}
		addon.SetName(name.Name)
		if err := r.Create(context.Background(), addon); err != nil && !apierrors.IsAlreadyExists(err) {
			r.log.Error(err, "couldn't generate addon", "addon", name, "component", obj.Meta.GetSelfLink())
		}
	}

	return
}
