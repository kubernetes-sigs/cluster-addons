package main

import (
	"bytes"
	"strings"
	"context"

	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
	"k8s.io/apimachinery/pkg/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/rbac/v1"
	"sigs.k8s.io/yaml"
)

func ParseYAMLtoRole(manifestStr string) (string, error){
	ctx := context.Background()
	objs, err := manifest.ParseObjects(ctx, manifestStr)
	if err != nil {
		return "", err
	}

	clusterRole := v1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: *name,
			Namespace: *ns,
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
	}
	// to deal with duplicates, we keep a map of all the kinds that has been added so far
	kindMap := make(map[string]bool)

	for _, obj := range objs.Items{
		// The generated role needs the rules from any role or clusterrole
		if obj.Kind == "Role" || obj.Kind == "ClusterRole" {
			unstruct := obj.UnstructuredObject()
			newClusterRole := v1.ClusterRole{}

			// Converting from unstructured to v1.ClusterRole
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstruct.Object, &newClusterRole)
			if err != nil {
				return "", err
			}
			clusterRole.Rules = append(clusterRole.Rules, newClusterRole.Rules...)
		}

		if !kindMap[obj.Group + "::" + obj.Kind] {
			newRule := v1.PolicyRule{
				APIGroups: []string{obj.Group},
				// needs plural of kind
				Resources: []string{resourceFromKind(obj.Kind)},
				Verbs: []string{"create", "update", "delete", "get"},
			}
			clusterRole.Rules = append(clusterRole.Rules, newRule)
			kindMap[obj.Group + "::" + obj.Kind] = true
		}
	}

	output, err := yaml.Marshal(&clusterRole)
	buf := bytes.NewBuffer(output)

	// if saName is passed in, generate YAML for rolebinding
	if *saName != "" {
		clusterRoleBinding := v1.ClusterRoleBinding{
			TypeMeta:   metav1.TypeMeta{
				Kind: "ClusterRoleBinding",
				APIVersion: "rbac.authorization.k8s.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: *ns,
				Name: *name + "-binding",
			},
			Subjects: []v1.Subject{
				{
					Kind: "ServiceAccount",
					Name: *saName,
					Namespace: *ns,
				},
			},
			RoleRef:    v1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind: "ClusterRole",
				Name: *name,
			},
		}

		outputBinding, err := yaml.Marshal(&clusterRoleBinding)
		if err != nil{
			return "", err
		}
		buf.WriteString("\n---\n\n")
		buf.Write(outputBinding)
	}
	return buf.String(), err
}

func resourceFromKind(kind string)  string{
	if string(kind[len(kind)-1]) == "s" {
		return strings.ToLower(kind) + "es"
	}
	if string(kind[len(kind)-1]) == "y" {
		return strings.ToLower(kind)[:len(kind) -1] + "ies"
	}
	return strings.ToLower(kind) + "s"
}
