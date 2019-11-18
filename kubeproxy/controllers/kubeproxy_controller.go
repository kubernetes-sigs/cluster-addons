package controllers

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/status"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "addon-operators/kubeproxy/api/v1alpha1"
)

var _ reconcile.Reconciler = &KubeProxyReconciler{}

// KubeProxyReconciler reconciles a KubeProxy object
type KubeProxyReconciler struct {
	client.Client
	Log logr.Logger

	declarative.Reconciler
}

// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=kubeproxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=addons.x-k8s.io,resources=kubeproxies/status,verbs=get;update;patch

func (r *KubeProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	addon.Init()

	labels := map[string]string{
		"k8s-app": "kubeproxy",
	}

	watchLabels := declarative.SourceLabel(mgr.GetScheme())

	if err := r.Reconciler.Init(mgr, &api.KubeProxy{},
		declarative.WithRawManifestOperation(injectFlags),
		declarative.WithRawManifestOperation(replaceNamespacePattern("{{.Namespace}}")),
		declarative.WithObjectTransform(declarative.AddLabels(labels)),
		declarative.WithObjectTransform(OverrideApiserver),
		declarative.WithOwner(declarative.SourceAsOwner),
		declarative.WithLabels(watchLabels),
		declarative.WithStatus(status.NewBasic(mgr.GetClient())),
		declarative.WithObjectTransform(addon.TransformApplicationFromStatus),
		declarative.WithManagedApplication(watchLabels),
		declarative.WithObjectTransform(addon.ApplyPatches),
	); err != nil {
		return err
	}

	c, err := controller.New("kubeproxy-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to KubeProxy
	err = c.Watch(&source.Kind{Type: &api.KubeProxy{}}, &handler.EnqueueRequestForObject{})
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

// replaceNamespacePattern fills in the namespace placeholder patterns with the actual namespace from the crd
func replaceNamespacePattern(nspatterns ...string) declarative.ManifestOperation {
	return func(ctx context.Context, o declarative.DeclarativeObject, manifest string) (string, error) {
		for _, pattern := range nspatterns {
			if strings.Index(manifest, pattern) != -1 {
				manifest = strings.Replace(manifest, pattern, o.GetNamespace(), -1)
			}
		}
		return manifest, nil
	}
}

func injectFlags(ctx context.Context, object declarative.DeclarativeObject, s string) (string, error) {
	o := object.(*api.KubeProxy)
	params := []string{
		"--v=2",
		"--iptables-sync-period=1m",
		"--iptables-min-sync-period=10s",
		"--ipvs-sync-period=1m",
		"--ipvs-min-sync-period=10s"}
	if o.Spec.ClusterCIDR != "" {
		params = append(params, "--cluster-cidr="+o.Spec.ClusterCIDR)
	}
	s = strings.Replace(s, "{{params}}", strings.Join(params, " "), -1)
	return s, nil
}

// SetEnvironmentVariables sets the env values on a pod template in the specified object
func SetEnvironmentVariables(o *manifest.Object, env map[string]string) error {
	// Using the unstructured library avoids problems when the manifest is using a newer API type thanwe are compiled with - as long as the field structure hasn't changed.
	if err := o.MutateContainers(func(container map[string]interface{}) error {
		envList, _, err := unstructured.NestedSlice(container, "env")
		if err != nil {
			return fmt.Errorf("error reading container env: %v", err)
		}
		for k, v := range env {
			foundK := false
			for _, e := range envList {
				m, ok := e.(map[string]interface{})
				if !ok {
					return fmt.Errorf("env var was not an object: %v", err)
				}
				name, found, err := unstructured.NestedString(m, "name")
				if err != nil {
					return err
				}
				if found && name == k {
					if err := unstructured.SetNestedField(m, v, "value"); err != nil {
						return err
					}
					foundK = true
				}
			}
			if !foundK {
				envList = append(envList, map[string]interface{}{
					"name":  k,
					"value": v,
				})
			}
		}
		// Sort env values by name so we have a consistent order
		sort.Slice(envList, func(i, j int) bool {
			mapI, okI := envList[i].(map[string]interface{})
			mapJ, okJ := envList[j].(map[string]interface{})
			if !okJ {
				return false
			}
			if !okI {
				return true
			}
			kI, foundI := mapI["name"]
			kJ, foundJ := mapJ["name"]
			if !foundJ {
				return false
			}
			if !foundI {
				return true
			}
			return kI.(string) < kJ.(string)
		})
		if err := unstructured.SetNestedSlice(container, envList, "env"); err != nil {
			return fmt.Errorf("error setting env vars: %v", err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// OverrideApiserver sets environment variables, hostNetwork & dns policy so that we don't rely on kube-proxy to reach the api server
func OverrideApiserver(ctx context.Context, o declarative.DeclarativeObject, manifest *manifest.Objects) error {
	for _, o := range manifest.Items {
		if o.Kind == "DaemonSet" {
			// KubeProxy (and kubelet) are special: it has to find the apiserver directly,
			// other clients use the VIP that kubeproxy configures.
			//
			// Because of this we need a special env var to allow for node kube-proxies to have a different endpoint than the operator
			// (If the operator runs on the master, it can use 127.0.0.1; we definitely can't use that for the node kube-proxy)
			// We still fall-back to the existing mechanisms - KUBERNETES_SERVICE_HOST, then a hard-coded default value
			master := os.Getenv("KUBEPROXY_KUBERNETES_SERVICE_HOST")
			if master == "" {
				master = os.Getenv("KUBERNETES_SERVICE_HOST")
			}
			if master == "" {
				master = "kubernetes-master"
				klog.Warningf("using fallback for KUBERNETES_SERVICE_HOST: %v", master)
			}
			port := os.Getenv("KUBERNETES_SERVICE_PORT")
			if port == "" {
				port = "443"
			}
			env := map[string]string{"KUBERNETES_SERVICE_HOST": master, "KUBERNETES_SERVICE_PORT": port}
			if err := SetEnvironmentVariables(o, env); err != nil {
				return err
			}
			if err := o.SetNestedField(true, "spec", "template", "spec", "hostNetwork"); err != nil {
				return fmt.Errorf("error setting hostNetwork: %v", err)
			}
			// To resolve the kubernetes-master field
			if err := o.SetNestedField("Default", "spec", "template", "spec", "dnsPolicy"); err != nil {
				return fmt.Errorf("error setting dnsPolicy: %v", err)
			}
		}
	}
	return nil
}
