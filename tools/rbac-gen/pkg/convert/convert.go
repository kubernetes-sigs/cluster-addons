package convert

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
)

type BuildRoleOptions struct {
	Name               string
	Namespace          string
	ServiceAccountName string
	Supervisory        bool

	// CRD is the name of the CRD to generate permissions for.
	CRD string

	// LimitResourceNames specifies that RBAC permissions should restrict to resource names in the manifest.
	LimitResourceNames bool

	// LimitNamespaces specifies that RBAC permissions should restrict to resource names in the manifest.
	LimitNamespaces bool

	// Format specifies the format we should write in (yaml or kubebuilder)
	Format string
}

// BuildRole parses the manifest and generates Role/ClusterRole objects for manipulating them.
func BuildRole(ctx context.Context, manifestStr string, opt BuildRoleOptions) ([]runtime.Object, error) {
	var objects []runtime.Object

	objs, err := manifest.ParseObjects(ctx, manifestStr)
	if err != nil {
		return nil, err
	}
	if len(objs.Blobs) != 0 {
		return nil, fmt.Errorf("unable to parse manifest fully")
	}

	// Build rules, keyed by namespace.  "" is used for cluster-scoped rules.
	ruleMap := make(map[string]*ruleSet)
	ruleMap[""] = &ruleSet{}

	// to deal with duplicates, we keep a map of all the kinds that has been added so far
	kindMap := make(map[string]bool)

	for _, obj := range objs.Items {
		ruleSetKey := ""
		if opt.LimitNamespaces {
			ruleSetKey = obj.Namespace
		}

		rules := ruleMap[ruleSetKey]
		if rules == nil {
			rules = &ruleSet{}
			ruleMap[ruleSetKey] = rules
		}

		// The generated role needs the rules from any role or clusterrole
		if obj.Kind == "Role" || obj.Kind == "ClusterRole" {
			if opt.Supervisory {
				continue
			}
			unstruct := obj.UnstructuredObject()
			newClusterRole := v1.ClusterRole{}

			// Converting from unstructured to v1.ClusterRole
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstruct.Object, &newClusterRole)
			if err != nil {
				return nil, err
			}
			rules.Add(newClusterRole.Rules...)
		}

		// needs plural of kind
		resource := ResourceFromKind(obj.Kind)

		if opt.LimitResourceNames {
			rules.Add(v1.PolicyRule{
				APIGroups:     []string{obj.Group},
				Resources:     []string{resource},
				ResourceNames: []string{obj.Name},
				Verbs:         []string{"update", "delete", "patch"},
			})

			rules.Add(v1.PolicyRule{
				APIGroups: []string{obj.Group},
				Resources: []string{resource},
				Verbs:     []string{"create"},
			})

			rules.Add(v1.PolicyRule{
				APIGroups: []string{obj.Group},
				Resources: []string{resource},
				Verbs:     []string{"get", "list", "watch"},
			})

		} else if !kindMap[obj.Group+"::"+obj.Kind] {
			newRule := v1.PolicyRule{
				APIGroups: []string{obj.Group},
				Resources: []string{resource},
				Verbs:     []string{"create", "update", "delete", "get"},
			}
			rules.Add(newRule)
			kindMap[obj.Group+"::"+obj.Kind] = true
		}
	}

	if opt.CRD != "" {
		crdGroupResource := schema.ParseGroupResource(opt.CRD)

		// TODO: Should we assume namespace scoped?
		rules := ruleMap[""]

		rules.Add(v1.PolicyRule{
			APIGroups: []string{crdGroupResource.Group},
			Resources: []string{crdGroupResource.Resource},
			Verbs:     []string{"get", "list", "patch", "update", "watch"},
		})
		rules.Add(v1.PolicyRule{
			APIGroups: []string{crdGroupResource.Group},
			Resources: []string{crdGroupResource.Resource + "/status"},
			Verbs:     []string{"get", "patch", "update"},
		})
	}

	// Normalize and wrap the rules in Role or ClusterRole objects
	for ns, rules := range ruleMap {
		rules.normalize()

		if ns != "" {
			role := &v1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:      opt.Name,
					Namespace: ns,
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "Role",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				Rules: rules.rules,
			}
			objects = append(objects, role)

			// if saName is passed in, generate YAML for rolebinding
			if opt.ServiceAccountName != "" {
				roleBinding := v1.RoleBinding{
					TypeMeta: metav1.TypeMeta{
						Kind:       "RoleBinding",
						APIVersion: "rbac.authorization.k8s.io/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      opt.Name + "-binding",
						Namespace: ns,
					},
					Subjects: []v1.Subject{
						{
							Kind:      "ServiceAccount",
							Name:      opt.ServiceAccountName,
							Namespace: opt.Namespace,
						},
					},
					RoleRef: v1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "Role",
						Name:     opt.Name,
					},
				}

				objects = append(objects, &roleBinding)
			}
		} else {
			role := &v1.ClusterRole{
				ObjectMeta: metav1.ObjectMeta{
					Name: opt.Name,
				},
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				Rules: rules.rules,
			}
			objects = append(objects, role)

			// if saName is passed in, generate YAML for rolebinding
			if opt.ServiceAccountName != "" {
				roleBinding := v1.ClusterRoleBinding{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterRoleBinding",
						APIVersion: "rbac.authorization.k8s.io/v1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: opt.Name + "-binding",
					},
					Subjects: []v1.Subject{
						{
							Kind:      "ServiceAccount",
							Name:      opt.ServiceAccountName,
							Namespace: opt.Namespace,
						},
					},
					RoleRef: v1.RoleRef{
						APIGroup: "rbac.authorization.k8s.io",
						Kind:     "ClusterRole",
						Name:     opt.Name,
					},
				}

				objects = append(objects, &roleBinding)
			}
		}
	}

	return objects, err
}

// ResourceFromKind returns the resource name for the given kind.
// Because we don't require a cluster, it assumes the kind -> resource mapping follows normal conventions.
func ResourceFromKind(kind string) string {
	if string(kind[len(kind)-1]) == "s" {
		return strings.ToLower(kind) + "es"
	}
	if string(kind[len(kind)-1]) == "y" {
		return strings.ToLower(kind)[:len(kind)-1] + "ies"
	}
	return strings.ToLower(kind) + "s"
}
