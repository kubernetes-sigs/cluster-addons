package xform

import (
	"context"
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Runnable defines the contract for transform objects.
type Runnable interface {
	Run(ctx context.Context, resourceList *framework.ResourceList) error
}

// TransformFunc defines the signature of a transformation function.
type TransformFunc func(ctx context.Context, resourceList *framework.ResourceList) error

// RunXform executes the specifies function against stdin/stdout.
func RunXform(ctx context.Context, fn TransformFunc) error {
	io := kio.ByteReadWriter{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}

	items, err := io.Read()
	if err != nil {
		return fmt.Errorf("failed to parse yaml: %w", err)
	}

	resourceList := &framework.ResourceList{
		Items: items,
	}

	if err := fn(ctx, resourceList); err != nil {
		return err
	}

	if err := io.Write(resourceList.Items); err != nil {
		return err
	}
	return nil
}

// RunFilters executes the specified filters against the provided resources.
func RunFilters(ctx context.Context, resourceList *framework.ResourceList, filters ...yaml.Filter) error {
	var out []*yaml.RNode
	for _, obj := range resourceList.Items {
		_, err := obj.Pipe(filters...)
		if err != nil {
			return err
		}
		out = append(out, obj)
	}

	resourceList.Items = out

	return nil
}
