package visitor

import (
	"fmt"

	"k8s.io/klog/v2"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func VisitResourceList(resourceList *framework.ResourceList, visitor visitor) error {
	ctx := &Context{}

	ctx.EnqueueVisitor(visitor)

	return ctx.visitResourceList(resourceList)
}

func visitNode(ctx *Context, path Path, node *yaml.Node, visitor visitor) error {
	switch node.Kind {
	case yaml.ScalarNode:
		return visitor.VisitScalar(ctx, path, node)

	case yaml.SequenceNode:
		if err := visitor.VisitSequence(ctx, path, node); err != nil {
			return err
		}
		n := len(node.Content)
		for i := 0; i < n; i += 2 {
			v := node.Content[i]
			childPath := path + "[]"
			if err := visitNode(ctx, childPath, v, visitor); err != nil {
				return err
			}
		}
		return nil

	case yaml.MappingNode:
		if err := visitor.VisitMap(ctx, path, &KubeMap{node: node}); err != nil {
			return err
		}
		n := len(node.Content)
		if n%2 != 0 {
			return fmt.Errorf("unexpected content length in MappingNode %v", path)
		}
		for i := 0; i < n; i += 2 {
			k := node.Content[i]
			ks, ok := AsString(k)
			if !ok {
				klog.Warningf("ignorning non-string MappingNode key at %v %v", path, k)
				continue
			}
			childPath := string(path) + "." + ks
			v := node.Content[i+1]
			if err := visitNode(ctx, Path(childPath), v, visitor); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unhandled yaml node kind %v", node.Kind)
	}
}

func AsString(n *yaml.Node) (string, bool) {
	if n.Kind != yaml.ScalarNode {
		return "", false
	}
	if n.Tag == "!!str" || n.Tag == "" {
		return n.Value, true
	}
	klog.Infof("Tag: %v", n.Tag)
	klog.Infof("Tag: %#v", n)
	return "", false
}
