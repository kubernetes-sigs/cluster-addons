package controllers

import (
	"context"
	"strings"

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

	api "sigs.k8s.io/cluster-addons/localnodedns/api/v1alpha1"
)

var _ reconcile.Reconciler = &LocalNodeDNSReconciler{}

// LocalNodeDNSReconciler reconciles a LocalNodeDNS object
type LocalNodeDNSReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	declarative.Reconciler
}

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=localnodedns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=localnodedns/status,verbs=get;update;patch

func (r *LocalNodeDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addon.Init()

	labels := map[string]string{
		"k8s-app": "localnodedns",
	}

	watchLabels := declarative.SourceLabel(mgr.GetScheme())

	if err := r.Reconciler.Init(mgr, &api.LocalNodeDNS{},
		declarative.WithRawManifestOperation(replaceVariables),
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(watchLabels),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		declarative.WithObjectTransform(addon.TransformApplicationFromStatus),
		declarative.WithManagedApplication(watchLabels),
		declarative.WithObjectTransform(addon.ApplyPatches),
	); err != nil {
		return err
	}

	c, err := controller.New("localnodedns-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to LocalNodeDNS
	err = c.Watch(&source.Kind{Type: &api.LocalNodeDNS{}}, &handler.EnqueueRequestForObject{})
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

func replaceVariables(ctx context.Context, object declarative.DeclarativeObject, s string) (string, error) {
	o := object.(*api.LocalNodeDNS)

	if o.Spec.DNSDomain == "" {
		o.Spec.DNSDomain = "cluster.local"
	}

	if o.Spec.DNSIP == "" {
		o.Spec.DNSIP = "169.254.20.10"
	}

	if o.Spec.ClusterIP == "" {
		o.Spec.ClusterIP = "10.96.0.10"
	}

	s = strings.Replace(s, "__PILLAR__LOCAL__DNS__", o.Spec.DNSIP, -1)
	s = strings.Replace(s, "__PILLAR__DNS__SERVER__", o.Spec.ClusterIP, -1)
	s = strings.Replace(s, "__PILLAR__DNS__DOMAIN__", o.Spec.DNSDomain, -1)

	return s, nil
}
