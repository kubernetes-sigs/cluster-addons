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

package coredns

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "sigs.k8s.io/addon-operators/coredns/pkg/apis/addons/v1alpha1"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/test/golden"
)

func TestCoreDNS(t *testing.T) {
	v := golden.NewValidator(t, api.SchemeBuilder)
	m := v.Manager()
	r := newReconciler(m)

	fakeClient := m.GetClient()
	objectmeta := metav1.ObjectMeta{Name: "kubernetes", Namespace: "default"}
	obj := &corev1.Service{ObjectMeta: objectmeta, Spec: corev1.ServiceSpec{ClusterIP: "10.96.0.1"}}
	fakeClient.Create(nil, obj)

	v.Validate(r.Reconciler)
}
