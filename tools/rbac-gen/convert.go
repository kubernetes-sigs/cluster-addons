package main

import (
	"strings"
	"context"

	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
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
	// to deal with duplicates, we keep a map of all the kinds that has been addeed so far
	kindMap := make(map[string]string)

	for _, obj := range objs.Items{
		// The generated role needs the rules from any role or clusterrole
		if obj.Kind == "Role" || obj.Kind == "ClusterRole" {
			unstruct := obj.UnstructuredObject()
			rules := unstruct.Object["rules"]
			//fmt.Println(rules)
			for _, rule := range rules.([]interface{}){
				rule := rule.(map[string]interface{})

				// we have to convert []interface{} to []string
				verbs := []string{}
				for _, intf := range rule["verbs"].([]interface{}) {
					verb := intf.(string)
					verbs = append(verbs, verb)
				}

				resources := []string{}
				for _, intf := range rule["resources"].([]interface{}) {
					resource := intf.(string)
					resources = append(resources, resource)
				}

				apiGroups := []string{}
				for _, intf := range rule["apiGroups"].([]interface{}) {
					apiGroup := intf.(string)
					apiGroups = append(apiGroups, apiGroup)
				}

				// TODO: Check for duplicates
				newRule := v1.PolicyRule{
					Verbs:           verbs,
					APIGroups:       apiGroups,
					Resources:       resources,
				}

				clusterRole.Rules = append(clusterRole.Rules, newRule)
			}
		}

		if _, ok := kindMap[obj.Kind]; !ok {
			newRule := v1.PolicyRule{
				APIGroups: []string{obj.Group},
				// needs plural of kind
				Resources: []string{resourceFromKind(obj.Kind)},
				Verbs: []string{"create", "update", "delete", "get"},
			}
			clusterRole.Rules = append(clusterRole.Rules, newRule)
			kindMap[obj.Kind] = ""
		}
	}

	output, err := yaml.Marshal(&clusterRole)
	return string(output), err
}

func resourceFromKind(kind string)  string{
	//map of apiresources that follow a different role
	if string(kind[len(kind)-1]) == "s" {
		return strings.ToLower(kind) + "es"
	}
	if string(kind[len(kind)-1]) == "y" {
		return strings.ToLower(kind)[:len(kind) -1] + "ies"
	}
	return strings.ToLower(kind) + "s"
}
