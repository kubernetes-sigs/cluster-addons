package convert

import (
	"fmt"
	"sort"
	"strings"

	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

// KubebuilderConverter converts Role/Cluster objects to kubebuilder directives.
// These can then be copied and pasted into the code.
type KubebuilderConverter struct {
	// Rules holds the generated kubebuilder rules
	Rules []string
}

// VisitObjects iterates over the provided Role/ClusterRule objects, generating equivalent kubebuilder statements.
func (c *KubebuilderConverter) VisitObjects(objects []runtime.Object) error {
	for _, obj := range objects {
		switch obj := obj.(type) {
		case *v1.ClusterRole:
			if err := c.visitRules(obj.Rules, ""); err != nil {
				return err
			}
		case *v1.Role:
			if err := c.visitRules(obj.Rules, obj.Namespace); err != nil {
				return err
			}
		case *v1.ClusterRoleBinding, *v1.RoleBinding:
			// Not kubebuilder
			klog.Infof("ignoring object of type %T for kubebuilder conversion", obj)
		default:
			return fmt.Errorf("unhandled type %T", obj)
		}
	}

	sort.Strings(c.Rules)

	return nil
}

// visitRules visits a set of policy rules, generating the equivalent kubebuilder rule
func (c *KubebuilderConverter) visitRules(rules []v1.PolicyRule, namespace string) error {
	for _, rule := range rules {
		def := "//+kubebuilder:rbac:"
		if len(rule.APIGroups) != 0 {
			def += "groups=" + strings.Join(rule.APIGroups, ";")
		}
		if namespace != "" {
			def += ",namespace=" + namespace
		}
		if len(rule.Resources) != 0 {
			def += ",resources=" + strings.Join(rule.Resources, ";")
		}
		if len(rule.ResourceNames) != 0 {
			def += ",resourceNames=" + strings.Join(rule.ResourceNames, ";")
		}
		if len(rule.Verbs) != 0 {
			def += ",verbs=" + strings.Join(rule.Verbs, ";")
		}
		if len(rule.NonResourceURLs) != 0 {
			def += ",urls=" + strings.Join(rule.NonResourceURLs, ";")
		}

		c.Rules = append(c.Rules, def)
	}

	return nil
}
