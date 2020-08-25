package controllers

import (
	"context"
	"log"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "sigs.k8s.io/cluster-addons/nodelocaldns/api/v1alpha1"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/test/golden"
)

func TestNodeLocalDNS(t *testing.T) {
	v := golden.NewValidator(t, api.SchemeBuilder)

	svc := defineResources()
	err := v.Manager().GetClient().Create(context.TODO(), svc)
	if err != nil {
		t.Fatalf("error creating Service: %v", err)
	}

	dr := &NodeLocalDNSReconciler{
		Client: v.Manager().GetClient(),
	}

	err = dr.setupReconciler(v.Manager())

	if err != nil {
		log.Fatal("Error creating reconciler: %V", err)
	}

	v.Validate(dr.Reconciler)
}

func defineResources() *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: metav1.NamespaceDefault, Name: "kubernetes"},
		Spec: corev1.ServiceSpec{
			ClusterIP: "169.254.20.1",
		},
	}

	return svc
}
