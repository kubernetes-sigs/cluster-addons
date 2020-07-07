package main

import (
	"flag"
	"fmt"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "~/.kube/config", "kubeconfig file")

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	clientset, err := kubernetes.NewForConfig(config)

	pod, err := clientset.CoreV1().Pods("").Get("kindnet-j4dpv", metav1.GetOptions{})
	args := os.Args[1:]

	if args[0] == "config" {
		fmt.Println(os.Getenv("KUBECONFOG"))
		os.Exit(0)
	}

	fmt.Println("I am a plugin named foo")
}
