package coredns

import (
	"k8s.io/apimachinery/pkg/util/intstr"

	addonsv1alpha1 "sigs.k8s.io/addon-operators/coredns-operator/pkg/apis/addons/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const DefaultCorefile = `
.:53 {
	errors
	health
	kubernetes cluster.local in-addr.arpa ip6.arpa {
	pods insecure
	upstream
	fallthrough in-addr.arpa ip6.arpa
	}
	prometheus :9153
	forward . /etc/resolv.conf
	cache 30
	loop
	reload
	loadbalance
}
`

func newLabelsForCR(cr *addonsv1alpha1.CoreDNS) map[string]string {
	return map[string]string{
		"k8s-app":                       "kube-dns",
		"kubernetes.io/cluster-service": "true",
		"kubernetes.io/name":            cr.Name,
	}
}

// newDeploymentForCR returns a Deployment with the same name/namespace as the cr
func newDeploymentForCR(cr *addonsv1alpha1.CoreDNS) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    newLabelsForCR(cr),
			Annotations: map[string]string{
				"seccomp.security.alpha.kubernetes.io/pod": "docker/default",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: newLabelsForCR(cr),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: newLabelsForCR(cr),
				},
				Spec: *newPodSpecForCR(cr),
			},
		},
	}
}

// newPodSpecForCR returns a PodSpec with the same name/namespace as the cr
func newPodSpecForCR(cr *addonsv1alpha1.CoreDNS) *corev1.PodSpec {
	// needed for pointers to constants in SecurityContext
	boolFalse := false
	boolTrue := true

	podspec := &corev1.PodSpec{
		ServiceAccountName: cr.Name,
		DNSPolicy:          corev1.DNSDefault,
		Volumes: []corev1.Volume{
			corev1.Volume{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: cr.Name,
						},
						Items: []corev1.KeyToPath{
							corev1.KeyToPath{
								Key:  "Corefile",
								Path: "Corefile",
							},
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:            "coredns",
				Image:           "k8s.gcr.io/coredns:1.3.1",
				ImagePullPolicy: corev1.PullIfNotPresent,
				Args: []string{
					"-conf",
					"/etc/coredns/Corefile",
				},
				Ports: []corev1.ContainerPort{
					corev1.ContainerPort{
						Name:          "dns-udp",
						ContainerPort: 53,
						Protocol:      corev1.ProtocolUDP,
					},
					corev1.ContainerPort{
						Name:          "dns-tcp",
						ContainerPort: 53,
						Protocol:      corev1.ProtocolTCP,
					},
					corev1.ContainerPort{
						Name:          "metrics",
						ContainerPort: 9153,
						Protocol:      corev1.ProtocolTCP,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					corev1.VolumeMount{
						Name:      "config",
						MountPath: "/etc/coredns",
						ReadOnly:  true,
					},
				},
				LivenessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/health",
							Port: intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 8080,
							},
							Scheme: corev1.URISchemeHTTP,
						},
					},
					InitialDelaySeconds: 60,
					TimeoutSeconds:      5,
					SuccessThreshold:    1,
					FailureThreshold:    5,
				},
				ReadinessProbe: &corev1.Probe{
					Handler: corev1.Handler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/health",
							Port: intstr.IntOrString{
								Type:   intstr.Int,
								IntVal: 8080,
							},
							Scheme: corev1.URISchemeHTTP,
						},
					},
				},
				SecurityContext: &corev1.SecurityContext{
					AllowPrivilegeEscalation: &boolFalse,
					Capabilities: &corev1.Capabilities{
						Add: []corev1.Capability{
							"NET_BIND_SERVICE",
						},
						Drop: []corev1.Capability{
							"all",
						},
					},
					ReadOnlyRootFilesystem: &boolTrue,
				},
			},
		},
	}

	if cr.Namespace == "kube-system" {
		podspec.PriorityClassName = "system-cluster-critical"
	}

	return podspec
}

// newServiceForCR returns a Service with the same name/namespace as the cr
func newServiceForCR(cr *addonsv1alpha1.CoreDNS) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    newLabelsForCR(cr),
		},
		Spec: corev1.ServiceSpec{
			Selector: newLabelsForCR(cr),
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Name:     "dns-udp",
					Port:     53,
					Protocol: corev1.ProtocolUDP,
				},
				corev1.ServicePort{
					Name:     "dns-tcp",
					Port:     53,
					Protocol: corev1.ProtocolTCP,
				},
				corev1.ServicePort{
					Name:     "metrics",
					Port:     9153,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

// newConfigMapForCR returns a ConfigMap with the same name/namespace as the cr
func newConfigMapForCR(cr *addonsv1alpha1.CoreDNS) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    newLabelsForCR(cr),
		},
		Data: map[string]string{
			"Corefile": cr.Spec.Corefile,
		},
	}
}

// newServiceAccountForCR returns a new ServiceAccount with the same name/namespace as the cr
func newServiceAccountForCR(cr *addonsv1alpha1.CoreDNS) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    newLabelsForCR(cr),
		},
	}
}

// newClusterRoleForCR returns a new ClusterRole with the same name as the cr
func newClusterRoleForCR(cr *addonsv1alpha1.CoreDNS) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.Namespace + ":" + cr.Name,
			Labels: newLabelsForCR(cr),
			Annotations: map[string]string{
				"rbac.authorization.kubernetes.io/autoupdate": "false", // apiserver does not manage this resource
			},
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"endpoints",
					"services",
					"pods",
					"namespaces",
				},
				Verbs: []string{
					"list",
					"watch",
				},
			},
			rbacv1.PolicyRule{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"nodes",
				},
				Verbs: []string{
					"get",
				},
			},
		},
	}
}

// newClusterRoleBindingForCR returns a new ClusterRoleBinding with the same name as the cr
func newClusterRoleBindingForCR(cr *addonsv1alpha1.CoreDNS) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.Namespace + ":" + cr.Name,
			Labels: newLabelsForCR(cr),
			Annotations: map[string]string{
				"rbac.authorization.kubernetes.io/autoupdate": "false", // apiserver does not manage this resource
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     cr.Namespace + ":" + cr.Name,
		},
		Subjects: []rbacv1.Subject{
			rbacv1.Subject{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      cr.Name,
				Namespace: cr.Namespace,
			},
		},
	}
}
