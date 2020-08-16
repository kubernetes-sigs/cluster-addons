module sigs.k8s.io/cluster-addons/generic

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v0.18.4
	sigs.k8s.io/cli-utils v0.16.0 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200721120314-d8f3ce551a4a
)

replace sigs.k8s.io/kubebuilder-declarative-pattern => github.com/SomtochiAma/kubebuilder-declarative-pattern v0.0.0-20200723151822-4aa1e9692ce6
