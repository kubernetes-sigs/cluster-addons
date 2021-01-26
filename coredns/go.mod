module sigs.k8s.io/cluster-addons/coredns

go 1.15

require (
	github.com/coredns/corefile-migration v1.0.10
	github.com/go-logr/logr v0.3.0
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.20.1
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.20.1
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.8.0
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20210322175944-13703a7722e0
)
