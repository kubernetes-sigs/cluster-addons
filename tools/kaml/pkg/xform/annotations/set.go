package annotations

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// BuildSetAnnotationsFilter builds a transform that sets the specified annotations.
func BuildSetAnnotationsFilter(annotations map[string]string) (yaml.Filter, error) {
	m := make(map[string]*yaml.Node)
	for k, v := range annotations {
		m[k] = &yaml.Node{Kind: yaml.ScalarNode, Value: v}
	}

	fieldPaths, err := xform.ParseFieldPaths(
		[]string{
			"metadata.annotations",
		})
	if err != nil {
		return nil, err
	}
	return &xform.FieldAdder{
		FieldPaths: fieldPaths,
		Add:        m,
	}, nil
}

// SetAnnotations describes a transform that sets annotations.
type SetAnnotations struct {
	Annotations map[string]string `json:"annotations"`
}

// Run applies the transform.
func (opt SetAnnotations) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	filter, err := BuildSetAnnotationsFilter(opt.Annotations)
	if err != nil {
		return err
	}

	return xform.RunFilters(ctx, resourceList, filter)
}

// AddSetAnnotationsCommand creates the cobra.Command.
func AddSetAnnotationsCommand(parent *cobra.Command) {
	var opt SetAnnotations

	cmd := &cobra.Command{
		Use:     "set-annotation",
		Aliases: []string{"set-annotations"},
		RunE: func(cmd *cobra.Command, args []string) error {
			m := make(map[string]string)
			for _, arg := range args {
				tokens := strings.SplitN(arg, "=", 2)
				m[tokens[0]] = tokens[1]
			}
			opt.Annotations = m
			return xform.RunXform(cmd.Context(), opt.Run)
		},
	}

	parent.AddCommand(cmd)
}
