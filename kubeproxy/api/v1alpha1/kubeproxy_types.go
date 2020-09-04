package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubeProxySpec defines the desired state of KubeProxy
type KubeProxySpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`

	ClusterCIDR string `json:"clusterCidr,omitempty"`
}

// KubeProxyStatus defines the observed state of KubeProxy
type KubeProxyStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// KubeProxy is the Schema for the  API
// +kubebuilder:subresource:status
type KubeProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeProxySpec   `json:"spec,omitempty"`
	Status KubeProxyStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &KubeProxy{}

func (o *KubeProxy) ComponentName() string {
	return "kubeproxy"
}

func (o *KubeProxy) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *KubeProxy) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *KubeProxy) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *KubeProxy) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true

// KubeProxyList contains a list of KubeProxy
type KubeProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeProxy{}, &KubeProxyList{})
}
