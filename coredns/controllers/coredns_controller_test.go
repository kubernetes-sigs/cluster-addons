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

package controllers

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "sigs.k8s.io/cluster-addons/coredns/api/v1alpha1"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/test/golden"
)

func TestCoreDNS(t *testing.T) {
	v := golden.NewValidator(t, api.SchemeBuilder)

	// Create the CoreDNS Resources to extract info if needed
	cm, svc, deploy := defineCoreDNSResources()
	err := v.Manager().GetClient().Create(context.TODO(), cm)
	if err != nil {
		t.Fatalf("error creating CoreDNS ConfigMap: %v", err)
	}
	err = v.Manager().GetClient().Create(context.TODO(), svc)
	if err != nil {
		t.Fatalf("error creating CoreDNS Service: %v", err)
	}
	err = v.Manager().GetClient().Create(context.TODO(), deploy)
	if err != nil {
		t.Fatalf("error creating CoreDNS Deployment: %v", err)
	}

	dr := &CoreDNSReconciler{
		Client: v.Manager().GetClient(),
	}
	err = dr.setupReconciler(v.Manager())
	if err != nil {
		t.Fatalf("creating reconciler: %v", err)
	}

	v.Validate(dr.Reconciler)
}

func defineCoreDNSResources() (*corev1.ConfigMap, *corev1.Service, *appsv1.Deployment) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coredns-xxxxxxxx",
			Namespace: metav1.NamespaceSystem,
		},
		Data: map[string]string{
			"Corefile": `.:53 {
    errors
    health {
        lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
	}

	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      coreDNSName,
			Namespace: metav1.NamespaceSystem,
			Labels: map[string]string{
				"k8s-app": "kube-dns",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Name: "coredns", Image: "k8s.gcr.io/coredns:1.7.0"}, // Note that this is used by the corefile upgrade logic
					},
					Volumes: []corev1.Volume{
						{Name: "config-volume", VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "coredns-xxxxxxxx"}}}},
					},
				},
			},
		},
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Namespace: metav1.NamespaceDefault, Name: "kubernetes"},
		Spec: corev1.ServiceSpec{
			ClusterIP: "10.96.0.1",
		},
	}

	return cm, svc, deploy
}
