package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CoreDNSSpec defines the desired state of CoreDNS
// +k8s:openapi-gen=true
type CoreDNSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	addonv1alpha1.CommonSpec

	// Corefile is a string representation of the operated CoreDNS's ConfigMap.
	// This string is hashed so that it's possible to do RollingUpdates of CoreDNS.
	Corefile string `json:"corefile,omitempty"`
	// ClusterDNS determines whether the operated CoreDNS reserves the 10th address of the cluster's service subnet.
	// Enabling this option makes the resulting CoreDNS service the canonical DNS server of the cluster.
	ClusterDNS bool `json:"clusterDNS,omitempty"`
}

// CoreDNSStatus defines the observed state of CoreDNS
// +k8s:openapi-gen=true
type CoreDNSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	addonv1alpha1.CommonStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoreDNS is the Schema for the coredns API
// +k8s:openapi-gen=true
type CoreDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoreDNSSpec   `json:"spec,omitempty"`
	Status CoreDNSStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoreDNSList contains a list of CoreDNS
type CoreDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoreDNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CoreDNS{}, &CoreDNSList{})
}
