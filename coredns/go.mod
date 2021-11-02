module sigs.k8s.io/cluster-addons/coredns

go 1.16

require (
	github.com/coredns/corefile-migration v1.0.14
	github.com/pkg/errors v0.9.1
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/klog/v2 v2.8.0
	sigs.k8s.io/controller-runtime v0.9.2
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20210630174303-f77bb4933dfb
)
