module sigs.k8s.io/cluster-addons/tools/generic-addon

go 1.13

require (
	github.com/gobuffalo/flect v0.2.1
	k8s.io/apimachinery v0.18.4
	sigs.k8s.io/cluster-addons/tools/rbac-gen v0.0.0-00010101000000-000000000000
	sigs.k8s.io/kubebuilder-declarative-pattern v0.0.0-20200605153943-bafc349a4a84
	sigs.k8s.io/yaml v1.2.0

)

replace (
	sigs.k8s.io/cluster-addons/tools/rbac-gen => ../../../cluster-addons/tools/rbac-gen
	sigs.k8s.io/kubebuilder-declarative-pattern => github.com/SomtochiAma/kubebuilder-declarative-pattern v0.0.0-20200723151822-4aa1e9692ce6
)
