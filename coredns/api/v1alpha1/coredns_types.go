/*
Copyright 2020 The Kubernetes Authors.

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

	DNSDomain string `json:"dnsDomain,omitempty"`
	DNSIP     string `json:"dnsIP,omitempty"`
	Corefile  string `json:"corefile,omitempty"`
}

// CoreDNSStatus defines the observed state of CoreDNS
type CoreDNSStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// CoreDNS is the Schema for the coredns API
type CoreDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoreDNSSpec   `json:"spec,omitempty"`
	Status CoreDNSStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &CoreDNS{}

func (o *CoreDNS) ComponentName() string {
	return "coredns"
}

func (o *CoreDNS) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *CoreDNS) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *CoreDNS) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *CoreDNS) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true

// CoreDNSList contains a list of CoreDNS
type CoreDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CoreDNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CoreDNS{}, &CoreDNSList{})
}
