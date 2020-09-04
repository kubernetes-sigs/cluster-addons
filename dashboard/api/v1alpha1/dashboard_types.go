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

// DashboardSpec defines the desired state of Dashboard
type DashboardSpec struct {
	addonv1alpha1.CommonSpec `json:",inline"`
	addonv1alpha1.PatchSpec  `json:",inline"`
}

// DashboardStatus defines the observed state of Dashboard
type DashboardStatus struct {
	addonv1alpha1.CommonStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// Dashboard is the Schema for the dashboards API
type Dashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DashboardSpec   `json:"spec,omitempty"`
	Status DashboardStatus `json:"status,omitempty"`
}

var _ addonv1alpha1.CommonObject = &Dashboard{}
var _ addonv1alpha1.Patchable = &Dashboard{}

func (o *Dashboard) ComponentName() string {
	return "dashboard"
}

func (o *Dashboard) CommonSpec() addonv1alpha1.CommonSpec {
	return o.Spec.CommonSpec
}

func (o *Dashboard) PatchSpec() addonv1alpha1.PatchSpec {
	return o.Spec.PatchSpec
}

func (o *Dashboard) GetCommonStatus() addonv1alpha1.CommonStatus {
	return o.Status.CommonStatus
}

func (o *Dashboard) SetCommonStatus(s addonv1alpha1.CommonStatus) {
	o.Status.CommonStatus = s
}

// +kubebuilder:object:root=true

// DashboardList contains a list of Dashboard
// +kubebuilder:subresource:status
type DashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dashboard{}, &DashboardList{})
}
