package prefix

import (
	"context"
	"strings"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/visitor"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
)

// fixupRefs fixes references to objects after a rename.
type fixupRefs struct {
	visitor.Visitor

	OldName string
	NewName string
	Group   string
	Version string
	Kind    string
}

// Run applies the transform.
func (opt fixupRefs) Run(ctx context.Context, resourceList *framework.ResourceList) error {
	return visitor.VisitResourceList(resourceList, &opt)
}

func (opt fixupRefs) VisitMap(ctx *visitor.Context, path visitor.Path, m *visitor.KubeMap) error {
	if !strings.HasSuffix(string(path), "Ref") {
		return nil
	}

	vals, err := m.ExtractStringFields()
	if err != nil {
		return err
	}

	name := vals["name"]
	if name != opt.OldName {
		return nil
	}

	if vals["kind"] != opt.Kind {
		return nil
	}

	if vals["apiGroup"] != opt.Group {
		return nil
	}

	if err := m.Set("name", opt.NewName); err != nil {
		return nil
	}
	return nil
}
