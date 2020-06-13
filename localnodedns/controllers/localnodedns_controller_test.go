package controllers

import (
	"log"
	api "sigs.k8s.io/cluster-addons/localnodedns/api/v1alpha1"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/test/golden"
	"testing"
)

func TestLocalnodedns(t *testing.T) {
	v := golden.NewValidator(t, api.SchemeBuilder)
	dr := &LocalNodeDNSReconciler{
		Client:      v.Manager().GetClient(),
	}

	err := dr.setupReconciler(v.Manager())

	if err != nil {
		log.Fatal("Error creating reconciler: %V", err)
	}

	v.Validate(dr.Reconciler)
}