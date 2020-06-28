package controllers

import (
	"log"
	"testing"

	api "sigs.k8s.io/cluster-addons/nodelocaldns/api/v1alpha1"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/test/golden"
)

func TestNodeLocalDNS(t *testing.T) {
	v := golden.NewValidator(t, api.SchemeBuilder)
	dr := &NodeLocalDNSReconciler{
		Client: v.Manager().GetClient(),
	}

	err := dr.setupReconciler(v.Manager())

	if err != nil {
		log.Fatal("Error creating reconciler: %V", err)
	}

	v.Validate(dr.Reconciler)
}
