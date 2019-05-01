package controller

import (
	"sigs.k8s.io/addon-operators/coredns-operator/pkg/controller/coredns"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, coredns.Add)
}
