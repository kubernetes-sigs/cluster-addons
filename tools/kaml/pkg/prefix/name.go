package prefix

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/visitor"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
)

// NamePrefix describes a transform that adds a prefix to object names.
type NamePrefix struct {
	visitor.Visitor

	Prefix string
	Kinds  []string
}

// Run applies the transform.
func (opt NamePrefix) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	return visitor.VisitResourceList(resourceList, &opt)
}

func contains(haystack []string, s string) bool {
	for _, v := range haystack {
		if s == v {
			return true
		}
	}
	return false
}

func (opt NamePrefix) VisitKubeObject(ctx *visitor.Context, obj *visitor.KubeObject) error {
	kind := obj.GetKind()
	if len(opt.Kinds) > 0 && !contains(opt.Kinds, kind) {
		return nil
	}

	apiVersion := obj.GetAPIVersion()
	group := ""
	version := ""

	if strings.Contains(apiVersion, "/") {
		tokens := strings.Split(apiVersion, "/")
		if len(tokens) != 2 {
			return fmt.Errorf("unexpected apiVersion %q", apiVersion)
		}
		group = tokens[0]
		version = tokens[1]
	} else {
		version = apiVersion
	}

	oldName := obj.GetName()
	newName := opt.Prefix + oldName

	obj.SetName(newName)

	ctx.EnqueueVisitor(&fixupRefs{
		OldName: oldName,
		NewName: newName,
		Group:   group,
		Version: version,
		Kind:    kind,
	})
	return nil
}

// AddNamePrefixCommand creates the cobra.Command.
func AddNamePrefixCommand(parent *cobra.Command) {
	var opt NamePrefix

	cmd := &cobra.Command{
		Use: "prefix-name",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected <prefix> to be passed")
			}
			opt.Prefix = args[0]
			return xform.RunXform(cmd.Context(), opt.Run)
		},
	}

	cmd.Flags().StringSliceVar(&opt.Kinds, "kind", opt.Kinds, "pass to only prefix objects of the specified kind; repeat for multiple kinds")

	parent.AddCommand(cmd)
}
