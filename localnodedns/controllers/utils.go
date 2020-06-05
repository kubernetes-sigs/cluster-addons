package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func findKubeProxyMode(ctx context.Context, c client.Client) (string, error) {
	kubeProxyConfigMap := &corev1.ConfigMap{}
	id := client.ObjectKey{Namespace: metav1.NamespaceSystem, Name: "kube-proxy"}

	err := c.Get(ctx, id, kubeProxyConfigMap)
	mode := kubeProxyConfigMap.Data["mode"]

	return mode, err
}
