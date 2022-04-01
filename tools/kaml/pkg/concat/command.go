package concat

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// AddConcatCommand creates the cobra.Command.
func AddConcatCommand(parent *cobra.Command) {
	var opt ConcatOptions

	cmd := &cobra.Command{
		Use: "concat",
		RunE: func(cmd *cobra.Command, args []string) error {
			opt.Files = args
			return Run(cmd.Context(), opt)
		},
	}

	parent.AddCommand(cmd)
}

// ConcatOptions holds the options for a concatention.
type ConcatOptions struct {
	Files []string
}

// Run applies the transform.
func Run(ctx context.Context, opt ConcatOptions) error {
	resourceList := &framework.ResourceList{}

	for _, p := range opt.Files {
		items, err := readItemsFromFile(p)
		if err != nil {
			return err
		}

		resourceList.Items = append(resourceList.Items, items...)
	}

	io := kio.ByteWriter{
		Writer: os.Stdout,
	}
	if err := io.Write(resourceList.Items); err != nil {
		return err
	}
	return nil
}

func readItemsFromFile(p string) ([]*yaml.RNode, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to open %q: %w", p, err)
	}
	defer f.Close()

	io := kio.ByteReader{
		Reader: f,
	}

	items, err := io.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	return items, nil
}
