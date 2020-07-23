package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GenericSpec defines the desired state of Generic
type GenericSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ObjectKind ObjectKind `json:"objectKind"`
}

// GenericStatus defines the observed state of Generic
type GenericStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Generic is the Schema for the generics API
type Generic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GenericSpec   `json:"spec,omitempty"`
	Status GenericStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &Generic{}

func (o *Generic) ComponentName() string {
	return "generic" // <--------
}

func (o *Generic) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *Generic) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *Generic) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *Generic) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
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
	Channel string `json:"channel,omitempty"`
}

// // GenericOperatorSpec defines the desired state of GenericOperator
// type GenericOperatorSpec struct {
// 	addonv1alpha1.CommonSpec `json:",inline"`
// 	addonv1alpha1.PatchSpec  `json:",inline"`

// 	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
// 	// Important: Run "make" to regenerate code after modifying this file
// 	ObjectKind ObjectKind `json:"objectKind,inline"`
// }

// // GenericOperatorStatus defines the observed state of Generic
// type GenericOperatorStatus struct {
// 	addonv1alpha1.CommonStatus `json:",inline"`

// 	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
// 	// Important: Run "make" to regenerate code after modifying this file
// }

// // +kubebuilder:object:root=true

// // GenericOperator is the Schema for the generics API
// type GenericOperator struct {
// 	metav1.TypeMeta   `json:",inline"`
// 	metav1.ObjectMeta `json:"metadata,omitempty"`

// 	Spec   GenericOperatorSpec   `json:"spec,omitempty"`
// 	Status GenericOperatorStatus `json:"status,omitempty"`
// }

// var _ addonv1alpha1.CommonObject = &GenericOperator{}

// func (o *GenericOperator) ComponentName() string {
// 	return "genericop" // <--------
// }

// func (o *GenericOperator) CommonSpec() addonv1alpha1.CommonSpec {
// 	return o.Spec.CommonSpec
// }

// func (o *GenericOperator) PatchSpec() addonv1alpha1.PatchSpec {
// 	return o.Spec.PatchSpec
// }

// func (o *GenericOperator) GetCommonStatus() addonv1alpha1.CommonStatus {
// 	return o.Status.CommonStatus
// }

// func (o *GenericOperator) SetCommonStatus(s addonv1alpha1.CommonStatus) {
// 	o.Status.CommonStatus = s
// }

// // +kubebuilder:object:root=true

// // GenericOperatorList contains a list of GenericOperator
// type GenericOperatorList struct {
// 	metav1.TypeMeta `json:",inline"`
// 	metav1.ListMeta `json:"metadata,omitempty"`
// 	Items           []GenericOperator `json:"items"`
// }

func init() {
	SchemeBuilder.Register(&Generic{}, &GenericList{})
}
