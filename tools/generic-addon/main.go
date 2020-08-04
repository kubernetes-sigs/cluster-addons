package main

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-addons/tools/rbac-gen/pkg/convert"

	//"k8s.io/apiextensions-apiserver/pkg/apiserver"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/gobuffalo/flect"
	"sigs.k8s.io/kubebuilder-declarative-pattern/pkg/patterns/addon/pkg/loaders"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	//seperator := "\n---\n"
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("Format: generic-addon [name] [group] [channel]")
		return fmt.Errorf("Please pass in enough arguments")
	}

	kind := args[0]
	group := args[1]
	channel := args[2]
	lower := strings.ToLower(kind)
	plural := flect.Pluralize(lower)

	resourceGroup := plural + "." + group
	addonCRD := strings.Replace(CRD, "<KIND>", kind, -1)
	addonCRD = strings.Replace(addonCRD, "<GROUP>", group, -1)
	addonCRD = strings.Replace(addonCRD, "<RESOURCEGROUP>", resourceGroup, -1)
	addonCRD = strings.Replace(addonCRD, "<PLURAL>", plural, -1)
	addonCRD = strings.Replace(addonCRD, "<SINGULAR>", lower, -1)

	crdRBAC := strings.Replace(CRDRBAC, "<PLURAL>", plural, -1)
	crdRBAC = strings.Replace(crdRBAC, "<GROUP>", group, -1)

	sampleYAML := strings.Replace(SAMPLEYAML, "<KIND>", kind, -1)
	sampleYAML = strings.Replace(sampleYAML, "<GROUP>", group, -1)
	sampleYAML = strings.Replace(sampleYAML, "<SINGULAR>", lower, -1)

	genericYaml := strings.Replace(GENERICYAML, "<KIND>", kind, -1)
	genericYaml = strings.Replace(genericYaml, "<GROUP>", group, -1)
	genericYaml = strings.Replace(genericYaml, "<CHANNEL>", channel, -1)

	strMap, err := getManifestFromChannel(channel, sampleYAML)
	if err != nil {
		return err
	}

	var allRbac string
	for _, manifestFile := range strMap {
		rbac, err := convert.ParseYAMLtoRole(manifestFile, "main-manager-role", "kube-system", lower+"system", false)
		if err != nil {
			return fmt.Errorf("error getting rbac from manifest: %v", err)
		}

		allRbac = allRbac + rbac + "\n---\n"
	}

	fmt.Println(addonCRD + crdRBAC + sampleYAML + genericYaml + allRbac)

	return nil
}

func getManifestFromChannel(url, sampleYaml string) (map[string]string, error) {
	manifestLoader, err := loaders.NewManifestLoader(url)
	if err != nil {
		return nil, err
	}

	json, err := yaml.YAMLToJSON([]byte(sampleYaml))
	o, _, err := unstructured.UnstructuredJSONScheme.Decode(json, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error converting yaml to json: %v", err)
	}

	m, err := manifestLoader.ResolveManifest(context.Background(), o)
	if err != nil {
		return nil, fmt.Errorf("error resolving manifest from channel:%v", err)
	}

	return m, nil
}

const CRD = `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
  creationTimestamp: null
  name: <RESOURCEGROUP>
spec:
  group: <GROUP>
  names:
    kind: <KIND>
    listKind: <KIND>List
    plural: <PLURAL>
    singular: <SINGULAR>
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: <KIND> is the Schema for the <PLURAL> API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: <KIND>Spec defines the desired state of Dashboard
          properties:
            channel:
              description: 'Channel specifies a channel that can be used to resolve
                a specific addon, eg: stable It will be ignored if Version is specified'
              type: string
            patches:
              items:
                type: object
              type: array
            version:
              description: Version specifies the exact addon version to be deployed,
                eg 1.2.3 It should not be specified if Channel is specified
              type: string
          type: object
        status:
          description: <KIND>Status defines the observed state of <KIND>
          properties:
            errors:
              items:
                type: string
              type: array
            healthy:
              type: boolean
          required:
          - healthy
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []

---
`

const CRDRBAC = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: manager-role
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  - secrets
  - serviceaccounts
  - services
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - proxy
- apiGroups:
  - ""
  resources:
  - services/proxy
  verbs:
  - get
- apiGroups:
  - '*'
  resources:
  - '*'
  verbs:
  - list
- apiGroups:
  - <GROUP>
  resources:
  - <PLURAL>
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - <GROUP>
  resources:
  - <PLURAL>/status
  verbs:
  - get
  - patch
  - update
- apiGroups:
  - app.k8s.io
  resources:
  - applications
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apps
  - extensions
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - metrics.k8s.io
  resources:
  - nodes
  - pods
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - rbac.authorization.k8s.io
  resources:
  - clusterrolebindings
  - clusterroles
  - rolebindings
  - roles
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch

---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: system

---
`

const SAMPLEYAML = `
apiVersion: <GROUP>/v1alpha1
kind: <KIND>
metadata:
  name: <SINGULAR>-sample
  namespace: kube-system
spec:
  channel: stable

---
`

const GENERICYAML = `
apiVersion: addons.x-k8s.io/v1alpha1
kind: Generic
metadata:
  name: generic-sample
spec:
  objectKind:
    kind: <KIND>
    version: v1alpha1
    group: <GROUP>
  channel: <CHANNEL>

---
`
