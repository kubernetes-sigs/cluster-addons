package labels

import (
	"context"

	"github.com/spf13/cobra"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// BuildRemoveLabelFilter builds a transform that removes labels matching the specificied predicate
func BuildRemoveLabelFilter(predicate func(key string) bool) (yaml.Filter, error) {
	fieldPaths, err := xform.ParseFieldPaths(
		[]string{
			"metadata.labels",
			"spec.selector",
			"spec.selector.matchLabels",
			"spec.template.metadata.labels",
		})
	if err != nil {
		return nil, err
	}
	return &xform.FieldClearer{
		FieldPaths: fieldPaths,
		Predicate:  predicate,
	}, nil
}

// RemoveLabel describes a transform that removes labels.
type RemoveLabel struct {
	Labels []string `json:"labels"`
}

// Run applies the transform.
func (opt RemoveLabel) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	predicate := func(key string) bool {
		for _, label := range opt.Labels {
			if label == key {
				return true
			}
		}
		return false
	}

	filter, err := BuildRemoveLabelFilter(predicate)
	if err != nil {
		return err
	}

	return xform.RunFilters(ctx, resourceList, filter)
}

// AddRemoveLabelsCommand creates the cobra.Command.
func AddRemoveLabelsCommand(parent *cobra.Command) {
	var opt RemoveLabel

	cmd := &cobra.Command{
		Use:     "remove-label",
		Aliases: []string{"remove-labels"},
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.Labels = append(opt.Labels, args...)
			return xform.RunXform(cmd.Context(), opt.Run)
		},
	}

	parent.AddCommand(cmd)
}
