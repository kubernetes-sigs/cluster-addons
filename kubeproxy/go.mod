module addon-operators/kubeproxy

go 1.12

require (
	github.com/go-logr/logr v0.1.0
	github.com/imdario/mergo v0.3.7 // indirect
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200512162422-ce639cbf6d4c
)
