package coredns

import (
	"context"
	"fmt"
	"net"

	"github.com/go-logr/logr"
	addonsv1alpha1 "sigs.k8s.io/addon-operators/coredns-operator/pkg/apis/addons/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_coredns")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new CoreDNS Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCoreDNS{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("coredns-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CoreDNS
	err = c.Watch(
		&source.Kind{Type: &addonsv1alpha1.CoreDNS{}},
		&handler.EnqueueRequestForObject{},
	)
	if err != nil {
		return err
	}

	// Watch for secondary types
	watchTypes := []runtime.Object{
		&appsv1.Deployment{},
		&corev1.Service{},
		&corev1.ConfigMap{},
		&corev1.ServiceAccount{},
		&rbacv1.ClusterRole{},
		&rbacv1.ClusterRoleBinding{},
	}
	for _, obj := range watchTypes {
		err = c.Watch(
			&source.Kind{Type: obj},
			&handler.EnqueueRequestForOwner{
				IsController: true,
				OwnerType:    &addonsv1alpha1.CoreDNS{},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileCoreDNS{}

// ReconcileCoreDNS reconciles a CoreDNS object
type ReconcileCoreDNS struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CoreDNS object and makes changes based on the state read
// and what is in the CoreDNS.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCoreDNS) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling CoreDNS")

	// Fetch the CoreDNS instance
	instance := &addonsv1alpha1.CoreDNS{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.Spec.Corefile == "" {
		instance.Spec.Corefile = DefaultCorefile
		reqLogger.Info("Updating CoreDNS with default Corefile")
		err = r.client.Update(context.TODO(), instance)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	foundDeployment := &appsv1.Deployment{}
	err = r.createOrFetch(reqLogger, instance, newDeploymentForCR(instance), foundDeployment)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	service := newServiceForCR(instance)
	dnsIP, err := r.calculateDNSClusterIP()
	if err != nil {
		return reconcile.Result{}, err
	}
	if instance.Spec.ClusterDNS {
		service.Spec.ClusterIP = dnsIP
	}
	foundService := &corev1.Service{}
	err = r.createOrFetch(reqLogger, instance, service, foundService)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		// The object was found already existing in the cluster
		// Check that the Service's ClusterIP is correct -- if not, update it
		if instance.Spec.ClusterDNS {
			if foundService.Spec.ClusterIP != dnsIP {
				reqLogger.Info("Re-creating Service", foundService.Namespace, "/", foundService.Name, "with operator-managed ClusterIP:", dnsIP)
				// ClusterIP field is immutable -- delete it, next reconcile will re-create
				err = r.client.Delete(context.TODO(), foundService)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		} else {
			if foundService.Spec.ClusterIP == dnsIP {
				reqLogger.Info("Re-creating Service", foundService.Namespace, "/", foundService.Name, "-- removing operator-managed ClusterIP:", dnsIP)
				// ClusterIP field is immutable -- delete it, next reconcile will re-create
				err = r.client.Delete(context.TODO(), foundService)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
	}

	foundConfigMap := &corev1.ConfigMap{}
	err = r.createOrFetch(reqLogger, instance, newConfigMapForCR(instance), foundConfigMap)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	} else if err == nil {
		// The object was found already existing in the cluster
		// Check that the ConfigMap's Corefile is correct -- if not, update it
		foundCorefile, inData := foundConfigMap.Data["Corefile"]
		if !inData || foundCorefile != instance.Spec.Corefile {
			foundConfigMap.Data["Corefile"] = instance.Spec.Corefile
			reqLogger.Info("Updating ConfigMap", foundConfigMap.Namespace, "/", foundConfigMap.Name, "with new Corefile")
			err = r.client.Update(context.TODO(), foundConfigMap)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	foundServiceAccount := &corev1.ServiceAccount{}
	err = r.createOrFetch(reqLogger, instance, newServiceAccountForCR(instance), foundServiceAccount)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	foundClusterRole := &rbacv1.ClusterRole{}
	err = r.createOrFetch(reqLogger, instance, newClusterRoleForCR(instance), foundClusterRole)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	foundClusterRoleBinding := &rbacv1.ClusterRoleBinding{}
	err = r.createOrFetch(reqLogger, instance, newClusterRoleBindingForCR(instance), foundClusterRoleBinding)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	// Reconcile is successful
	return reconcile.Result{}, nil
}

func (r *ReconcileCoreDNS) calculateDNSClusterIP() (string, error) {
	kubernetesService := &corev1.Service{}
	id := client.ObjectKey{Namespace: "default", Name: "kubernetes"}
	if err := r.client.Get(context.TODO(), id, kubernetesService); err != nil {
		return "", fmt.Errorf("error getting service %s: %v", id, err)
	}

	ip := net.ParseIP(kubernetesService.Spec.ClusterIP)
	if ip == nil {
		return "", fmt.Errorf("cannot parse kubernetes ClusterIP %q", kubernetesService.Spec.ClusterIP)
	}

	// The kubernetes Service ClusterIP is the 1st IP in the Service Subnet.
	// Increment the right-most byte by 9 to get to the 10th address, canonically used for kube-dns.
	// This works for both IPV4, IPV6, and 16-byte IPV4 addresses.
	ip[len(ip)-1] += 9

	return ip.String(), nil
}

// createOrFetch creates an object or populates found with the matching object from the cluster.
// It returns a notFound error if the object is created.
func (r *ReconcileCoreDNS) createOrFetch(reqLogger logr.Logger, instance metav1.Object, obj, found runtime.Object) error {
	meta, ok := obj.(metav1.Object)
	if !ok {
		return fmt.Errorf("Meta conversion failed for obj: %+v", obj)
	}

	// Set CoreDNS instance as the object owner and controller
	if err := controllerutil.SetControllerReference(instance, meta, r.scheme); err != nil {
		return err
	}

	// Check if this Object already exists
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}
	err = r.client.Get(
		context.TODO(),
		key,
		found,
	)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new obj", "obj.Kind", obj.GetObjectKind().GroupVersionKind(), "obj.Namespace", meta.GetNamespace(), "obj.Name", meta.GetName())
		createErr := r.client.Create(context.TODO(), obj)
		if createErr != nil {
			return err
		}
		return err
	} else if err != nil {
		return err
	} else {
		// obj already exists - just log
		foundGVK := found.GetObjectKind().GroupVersionKind()
		foundMeta, ok := found.(metav1.Object)
		if !ok {
			return fmt.Errorf("Meta conversion failed for found: %+v", found)
		}
		reqLogger.Info("Skip reconcile: obj already exists", "obj.Kind", foundGVK, "obj.Namespace", foundMeta.GetNamespace(), "obj.Name", foundMeta.GetName())

		fmt.Println("\n", found.GetObjectKind().GroupVersionKind().Kind)
	}

	return nil
}
