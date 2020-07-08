package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var (
	kubeconfig =  flag.String("kubeconfig", "~/.kube/config", "kubeconfig file")
	ns =  flag.String("ns", "", "namespace")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

}

func run() error{
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil{
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil
	}

	listOpts := metav1.ListOptions{}
	pods, err := clientset.CoreV1().Pods(*ns).List(listOpts)
	if err != nil {
		return err
	}

	printPods(pods)
}

func printPods(podlist *v1.PodList) {
	template := "%-32s%-8s-%-8s\n"
	fmt.Printf(template, "NAMESPACE","NAME", "STATUS")

	for _, pod := range podlist.Items {
		fmt.Printf(template,
			pod.Namespace,
			pod.Name,
			pod.Status,
		)
	}
}
