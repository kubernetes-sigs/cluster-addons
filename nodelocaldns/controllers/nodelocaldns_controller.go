package controllers

import (
	"context"
	"strings"

	"k8s.io/klog"

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

	api "sigs.k8s.io/cluster-addons/nodelocaldns/api/v1alpha1"
)

var _ reconcile.Reconciler = &NodeLocalDNSReconciler{}

// NodeLocalDNSReconciler reconciles a NodeLocalDNS object
type NodeLocalDNSReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	declarative.Reconciler
	watchLabels declarative.LabelMaker
}

func (r *NodeLocalDNSReconciler) setupReconciler(mgr ctrl.Manager) error {
	addon.Init()

	labels := map[string]string{
		"k8s-app": "nodelocaldns",
	}

	r.watchLabels = declarative.SourceLabel(mgr.GetScheme())

	return r.Reconciler.Init(mgr, &api.NodeLocalDNS{},
		declarative.WithRawManifestOperation(replaceVariables(mgr)),
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(r.watchLabels),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		declarative.WithObjectTransform(addon.TransformApplicationFromStatus),
		declarative.WithManagedApplication(r.watchLabels),
		declarative.WithObjectTransform(addon.ApplyPatches),
	)
}

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=nodelocaldns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=nodelocaldns/status,verbs=get;update;patch

func (r *NodeLocalDNSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := r.setupReconciler(mgr); err != nil {
		return err
	}

	c, err := controller.New("nodelocaldns-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to NodeLocalDNS
	err = c.Watch(&source.Kind{Type: &api.NodeLocalDNS{}}, &handler.EnqueueRequestForObject{})
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

func replaceVariables(mgr ctrl.Manager) declarative.ManifestOperation {
	return func(ctx context.Context, object declarative.DeclarativeObject, s string) (string, error) {
		o := object.(*api.NodeLocalDNS)
		kubeProxyMode, err := findKubeProxyMode(ctx, mgr.GetClient())
		if err != nil {
			klog.Warningf("error determining kube-proxy mode, defaulting to iptables: %v", err)
		}

		// TODO: port findClusterIP and getDNSDomain from coredns/controllers/utils in the kubebuilder-declarative
		// -pattern repo and use it here
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
		s = strings.Replace(s, "__PILLAR__DNS__DOMAIN__", o.Spec.DNSDomain, -1)

		if kubeProxyMode == "ipvs" {
			s = strings.Replace(s, "__PILLAR__DNS__SERVER__", "", -1)
			s = strings.Replace(s, "__PILLAR__CLUSTER__DNS__", o.Spec.ClusterIP, -1)
		} else {
			s = strings.Replace(s, "__PILLAR__DNS__SERVER__", o.Spec.ClusterIP, -1)
		}

		return s, nil
	}
}
