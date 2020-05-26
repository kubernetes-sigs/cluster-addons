module sigs.k8s.io/cluster-addons/coredns

go 1.13

require (
	github.com/coredns/corefile-migration v1.0.8
	github.com/go-logr/logr v0.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/pkg/errors v0.8.1
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200512162422-ce639cbf6d4c
)
