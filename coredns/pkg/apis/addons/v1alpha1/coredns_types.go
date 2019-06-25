/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonv1alpha1 "sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/apis/v1alpha1"
)

// CoreDNSSpec defines the desired state of CoreDNS
type CoreDNSSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`
}

// CoreDNSStatus defines the observed state of CoreDNS
type CoreDNSStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`
}

var _ addonv1alpha1.CommonObject = &CoreDNS{}
var _ addonv1alpha1.Patchable = &CoreDNS{}

func (c *CoreDNS) ComponentName() string {
	return "coredns"
}

func (c *CoreDNS) CommonSpec() addonv1alpha1.CommonSpec {
	return c.Spec.CommonSpec
}

func (c *CoreDNS) GetCommonStatus() addonv1alpha1.CommonStatus {
	return c.Status.CommonStatus
}

func (c *CoreDNS) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	c.Status.CommonStatus = s
}

func (c *CoreDNS) PatchSpec() addonv1alpha1.PatchSpec {
	return c.Spec.PatchSpec
}

// +genclient
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
