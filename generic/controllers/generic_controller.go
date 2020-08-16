package controllers

import (
	"fmt"
	"strings"

	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/loaders"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/status"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//api "sigs.k8s.io/cluster-addons/generic/api/v1alpha1"
)

var _ reconcile.Reconciler = &GenericReconciler{}

// GenericReconciler reconciles a Generic object
type GenericReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	declarative.Reconciler
	GVK     schema.GroupVersionKind
	Channel string
}

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=generics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=generics/status,verbs=get;update;patch

func (r *GenericReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addon.Init()

	labels := map[string]string{
		"k8s-app": strings.ToLower(r.GVK.Kind),
	}

	watchLabels := declarative.SourceLabel(mgr.GetScheme())

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(r.GVK)

	mc, err := loaders.NewManifestLoader(r.Channel)
	if err != nil {
		return fmt.Errorf("unable to create manifest loader: %v", err)
	}

	if err := r.Reconciler.Init(mgr, u,
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(watchLabels),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		// TODO: add an application to your manifest:
		// declarative.WithObjectTransform(addon.TransformApplicationFromStatus),
		// TODO: add an application to your manifest:
		// declarative.WithManagedApplication(watchLabels),
		declarative.WithObjectTransform(addon.ApplyPatches),
		declarative.WithManifestController(mc),
	); err != nil {
		return err
	}

	c, err := controller.New(strings.ToLower(r.GVK.Kind)+"-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to objectKind
	err = c.Watch(&source.Kind{Type: u}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to deployed objects
	_, err = declarative.WatchAll(mgr.GetConfig(), c, r, watchLabels)
	if err != nil {
		return err
	}

	return nil
}
