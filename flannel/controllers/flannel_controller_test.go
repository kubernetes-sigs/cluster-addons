package controllers

import (
	"testing"

	api "sigs.k8s.io/cluster-addons/flannel/api/v1alpha1"

	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/test/golden"
)

func TestFlannel(t *testing.T) {
	v := golden.NewValidator(t, api.SchemeBuilder)

	dr := &FlannelReconciler{
		Client: v.Manager().GetClient(),
	}

	err := dr.setupReconciler(v.Manager())
	if err != nil {
		t.Fatalf("creating reconciler: %v", err)
	}

	v.Validate(dr.Reconciler)
}
