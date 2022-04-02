package visitor

import (
	"k8s.io/klog/v2"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Context struct {
	visitorQueue []visitor
}

func (c *Context) EnqueueVisitor(visitor visitor) {
	c.visitorQueue = append(c.visitorQueue, visitor)
}

func (ctx *Context) visitResourceList(resourceList *framework.ResourceList) error {
	for len(ctx.visitorQueue) > 0 {
		visitor := ctx.visitorQueue[0]
		ctx.visitorQueue = ctx.visitorQueue[1:]

		for _, item := range resourceList.Items {
			var kubeObject *KubeObject
			switch item.YNode().Kind {
			case yaml.MappingNode:
				m := &KubeMap{node: item.YNode()}
				kubeObject = &KubeObject{obj: m}
			default:
				klog.Warningf("unexpected top-level yaml kind %v", item.YNode().Kind)
			}

			if kubeObject != nil {
				if err := visitor.VisitKubeObject(ctx, kubeObject); err != nil {
					return err
				}
			}

			if err := visitNode(ctx, "", item.YNode(), visitor); err != nil {
				return err
			}
		}
	}
	return nil
}
