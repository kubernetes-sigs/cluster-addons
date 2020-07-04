package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if args[0] == "config" {
		fmt.Println(os.Getenv("KUBECONFOG"))
		os.Exit(0)
	}

	fmt.Println("I am a plugin named foo")
}
