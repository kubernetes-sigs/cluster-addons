package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/declarative/pkg/manifest"
	"sigs.k8s.io/yaml"

)

type roleStruct struct {
	Kind string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata metadata `json:"metadata"`
	Rules []*rule `json:"rules"`
}

type metadata struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
}

type rule struct {
	ApiGroups []string `json:"apiGroups"`
	Resources []string `json:"resources"`
	Verbs []string `json:"verbs"`
}

var (
	yamlFile = flag.String("yaml", "manifest.yaml", "yaml file from which the rbac will be generated.")
	name = flag.String("name", "generated-role", "name of role to be generated")
	ns = flag.String("ns", "kube-system", "namespace of the role to be generated")
	out = flag.String("out.yaml", "kube-system", "name of output file")
)

func main() {
	//	read yaml file passed in from cmd
	flag.Parse()

	bytes, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		log.Fatalf("Error reading files: %v", err)
	}
	// generate Group and Kind
	ctx := context.Background()
	objs, err := manifest.ParseObjects(ctx, string(bytes))

	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	roleInterface := roleStruct{
		ApiVersion: "rbac.authorization.k8s.io/v1",
		Kind: "Role",
		Metadata: metadata{
			Name: *name,
			Namespace: *ns,
		},
	}

	// to deal with duplicates, we keep a map of all the kinds that has been addeed so far
	m := make(map[string]string)

	for _, obj := range objs.Items{

		if _, ok := m[obj.Kind]; !ok {
			newRule := rule {
				ApiGroups: []string{obj.Group},
				// needs plural of kind
				Resources: []string{resourceFromKind(obj.Kind)},
				Verbs: []string{"create", "update", "delete", "get"},
			}
			roleInterface.Rules = append(roleInterface.Rules, &newRule)
			m[obj.Kind] = ""
		}

	}

	out, err := yaml.Marshal(&roleInterface)
	if err != nil {
		log.Fatalf("Error parsing yaml: %v", err)
	}

	//fmt.Println(string(out))
	err = ioutil.WriteFile("out.yaml", out, 0644)
	if err != nil {
		log.Fatal(err)
	}

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
