package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	ns = flag.String("ns", "", "namespace")
)

func main() {
	var kubeconfig *string
	if homeDir := homedir.HomeDir(); homeDir != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(homeDir, ".kube", "config"), "kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "kubeconfig file")
	}

	flag.Parse()

	if err := run(kubeconfig); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}

func run(kubeconfig *string) error {
	if len(os.Args) < 3 {
		fmt.Println("Please complete command: kubectl ownerref [KIND] [NAME]")
		os.Exit(1)
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return err
	}

	discoveryClientSet, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}

	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	apiResourceList, err := discoveryClientSet.ServerPreferredResources()
	if err != nil {
		return nil
	}

	for _, apiResource := range apiResourceList {
		for _, resource := range apiResource.APIResources {
			group, version := getGroupVersion(apiResource.GroupVersion)
			res := schema.GroupVersionResource{
				Group:    group,
				Version:  version,
				Resource: resource.Name,
			}
			PrintResource(clientset, res, os.Args[1], os.Args[2])
		}
	}

	return nil
}

func PrintResource(c dynamic.Interface, resource schema.GroupVersionResource, kind, name string) error {
	listOpts := metav1.ListOptions{}
	resourceList, err := c.Resource(resource).Namespace(*ns).List(context.Background(), listOpts)
	if err != nil {
		return err
	}

	PrintList(resourceList, kind, name)
	return nil
}

// PrintList is a function for printing all kinds of resources
func PrintList(list *unstructured.UnstructuredList, kind, name string) {
	template := "%-32s%-32s%-8s\n"
	if len(list.Items) == 0 {
		return
	}

	var ownerItems []unstructured.Unstructured
	for _, item := range list.Items {
		for _, ownerRef := range item.GetOwnerReferences() {
			if strings.ToLower(ownerRef.Kind) == kind && ownerRef.Name == name {
				ownerItems = append(ownerItems, item)
			}
		}
	}

	if len(ownerItems) == 0 {
		return
	}

	fmt.Printf(template, "KIND", "NAMESPACE", "NAME")
	for _, ownerItem := range ownerItems {
		fmt.Printf(template,
			ownerItem.GetKind(),
			ownerItem.GetNamespace(),
			ownerItem.GetName(),
		)
	}
	fmt.Println()
}

func getGroupVersion(s string) (string, string) {
	slice := strings.Split(s, "/")
	if len(slice) == 1 {
		return "", slice[0]
	}

	return slice[0], slice[1]
}
