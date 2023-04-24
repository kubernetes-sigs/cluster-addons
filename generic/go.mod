module sigs.k8s.io/cluster-addons/generic

go 1.16

require (
	github.com/go-logr/logr v0.1.0
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	k8s.io/apimachinery v0.18.4
	k8s.io/client-go v0.18.4
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200721120314-d8f3ce551a4a
)

replace sigs.k8s.io/kubebuilder-declarative-pattern => github.com/SomtochiAma/kubebuilder-declarative-pattern v0.0.0-20200723151822-4aa1e9692ce6
