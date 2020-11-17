package xform

import "sigs.k8s.io/kustomize/kyaml/yaml"

// FieldClearerFunc removes field or map keys that match the predicate.
type fieldClearerFunc struct {
	// Predicate matches against the name of the field or key in the map.
	Predicate func(string) bool
}

// FieldClearerFunc is a yaml.Filter
var _ yaml.Filter = &fieldClearerFunc{}

// Filter implements the yaml filtering logic.
func (c fieldClearerFunc) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	if err := yaml.ErrorIfInvalid(rn, yaml.MappingNode); err != nil {
		return nil, err
	}

	var keep []*yaml.Node

	content := rn.Content()
	for i := 0; i < len(content); i += 2 {
		// if name matches, remove these 2 elements from the list because
		// they are treated as a fieldName/fieldValue pair.
		if !c.Predicate(content[i].Value) {
			keep = append(keep, content[i])
			if len(content) > i+1 {
				keep = append(keep, content[i+1])
			}
		}
	}
	rn.YNode().Content = keep

	return nil, nil
}

// FieldClearer removes fields at the specified FieldPaths, that match the predicate.
type FieldClearer struct {
	FieldPaths []FieldPath

	// Predicate matches against the name of the field or key in the map.
	Predicate func(string) bool
}

// FieldClearer is a yaml.Filter
var _ yaml.Filter = &FieldClearer{}

// Filter implements the yaml filtering logic.
func (c FieldClearer) Filter(rn *yaml.RNode) (*yaml.RNode, error) {
	for _, fieldPath := range c.FieldPaths {
		_, err := rn.Pipe(
			yaml.PathGetter{Path: fieldPath},
			fieldClearerFunc{Predicate: c.Predicate})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
