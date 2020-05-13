package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LocalNodeDNSSpec defines the desired state of LocalNodeDNS
type LocalNodeDNSSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// LocalNodeDNSStatus defines the observed state of LocalNodeDNS
type LocalNodeDNSStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// LocalNodeDNS is the Schema for the localnodedns API
type LocalNodeDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LocalNodeDNSSpec   `json:"spec,omitempty"`
	Status LocalNodeDNSStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &LocalNodeDNS{}

func (o *LocalNodeDNS) ComponentName() string {
	return "localnodedns"
}

func (o *LocalNodeDNS) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *LocalNodeDNS) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *LocalNodeDNS) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *LocalNodeDNS) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true

// LocalNodeDNSList contains a list of LocalNodeDNS
type LocalNodeDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LocalNodeDNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LocalNodeDNS{}, &LocalNodeDNSList{})
}
