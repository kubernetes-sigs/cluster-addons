package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenericSpec defines the desired state of Generic
type GenericSpec struct {
	ObjectKind ObjectKind `json:"objectKind"`
	Channel    string     `json:"channel"`
}

// GenericStatus defines the observed state of Generic
type GenericStatus struct {
	// addonv1alpha1.CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// Generic is the Schema for the generics API
type Generic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GenericSpec   `json:"spec,omitempty"`
	Status GenericStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GenericList contains a list of Generic
type GenericList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Generic `json:"items"`
}

type ObjectKind struct {
	Kind    string `json:"kind"`
	Group   string `json:"group"`
	Version string `json:"version"`
}

func init() {
	SchemeBuilder.Register(&Generic{}, &GenericList{})
}
