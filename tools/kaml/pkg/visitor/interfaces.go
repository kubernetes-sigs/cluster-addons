package visitor

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type visitor interface {
	VisitSequence(ctx *Context, path Path, node *yaml.Node) error
	VisitMap(ctx *Context, path Path, node *KubeMap) error
	VisitScalar(ctx *Context, path Path, node *yaml.Node) error

	VisitKubeObject(ctx *Context, obj *KubeObject) error
}

type Visitor struct {
}

func (v *Visitor) VisitSequence(ctx *Context, path Path, node *yaml.Node) error {
	return nil
}

func (v *Visitor) VisitScalar(ctx *Context, path Path, node *yaml.Node) error {
	return nil
}

func (v *Visitor) VisitMap(ctx *Context, path Path, node *KubeMap) error {
	return nil
}

func (v *Visitor) VisitKubeObject(ctx *Context, obj *KubeObject) error {
	return nil
}
