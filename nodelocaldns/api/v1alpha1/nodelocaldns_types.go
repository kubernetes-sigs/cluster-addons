package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// nodelocaldnsSpec defines the desired state of nodelocaldns
type NodeLocalDNSSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DNSDomain string `json:"dnsDomain,omitempty"`
	DNSIP     string `json:"dnsIP,omitempty"`
	ClusterIP string `json:"clusterIP,omitempty"`
}

// nodelocaldnsStatus defines the observed state of nodelocaldns
type NodeLocalDNSStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`

	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// nodelocaldns is the Schema for the nodelocaldns API
type NodeLocalDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeLocalDNSSpec   `json:"spec,omitempty"`
	Status NodeLocalDNSStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &NodeLocalDNS{}

func (o *NodeLocalDNS) ComponentName() string {
	return "nodelocaldns"
}

func (o *NodeLocalDNS) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *NodeLocalDNS) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *NodeLocalDNS) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *NodeLocalDNS) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true

// NodelocaldnsList contains a list of nodelocaldns
type NodeLocalDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeLocalDNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeLocalDNS{}, &NodeLocalDNSList{})
}
