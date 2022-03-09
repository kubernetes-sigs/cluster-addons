package main

import (
	"context"
	goflags "flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-addons/tools/rbac-gen/pkg/convert"
)

func main() {
	err := run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	rootCommand := BuildGenerateCommand(ctx)

	fs := goflags.NewFlagSet("", goflags.PanicOnError)
	klog.InitFlags(fs)
	rootCommand.PersistentFlags().AddGoFlagSet(fs)

	rootCommand.SilenceErrors = true
	rootCommand.SilenceUsage = true

	if err := rootCommand.Execute(); err != nil {
		return err
	}
	return nil
}

func BuildGenerateCommand(ctx context.Context) *cobra.Command {
	yamlFile := "manifest.yaml"
	out := ""

	var opt convert.BuildRoleOptions
	opt.Name = "generated-role"
	opt.Namespace = "kube-system"
	opt.Format = "yaml"

	cmd := &cobra.Command{
		Use: "generate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGenerate(ctx, yamlFile, out, opt)
		},
	}

	cmd.Flags().StringVar(&yamlFile, "yaml", yamlFile, "yaml file from which the rbac will be generated.")
	cmd.Flags().StringVar(&opt.Name, "name", opt.Name, "name of role to be generated")
	cmd.Flags().StringVar(&opt.ServiceAccountName, "sa-name", opt.ServiceAccountName, "name of service account the role should be bound to")
	cmd.Flags().StringVar(&opt.Namespace, "ns", opt.Namespace, "namespace of the role to be generated")
	cmd.Flags().StringVar(&out, "out", out, "name of output file")
	cmd.Flags().BoolVar(&opt.Supervisory, "supervisory", opt.Supervisory, "outputs role for operator in supervisory mode")
	cmd.Flags().StringVar(&opt.CRD, "crd", opt.CRD, "CRD to generate")
	cmd.Flags().BoolVar(&opt.LimitResourceNames, "limit-resource-names", opt.LimitResourceNames, "Limit to resource names in the manifest")
	cmd.Flags().BoolVar(&opt.LimitNamespaces, "limit-namespaces", opt.LimitNamespaces, "Limit to namespaces in the manifest")

	cmd.Flags().StringVar(&opt.Format, "format", opt.Format, "Format to write in (yaml, kubebuilder)")

	return cmd
}

func RunGenerate(ctx context.Context, yamlFile string, out string, opt convert.BuildRoleOptions) error {
	//	read yaml from file or stdin
	in := ""
	if yamlFile == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		in = string(b)
	} else {
		b, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			return err
		}
		in = string(b)
	}

	// build roles for objects in yaml
	objects, err := convert.BuildRole(ctx, in, opt)
	if err != nil {
		return err
	}

	var output []byte
	switch opt.Format {
	case "yaml":
		y, err := convert.ToYAML(objects)
		if err != nil {
			return err
		}
		output = y

	case "kubebuilder":
		// convert to kubebuilder format and output
		var conv convert.KubebuilderConverter
		if err := conv.VisitObjects(objects); err != nil {
			return err
		}
		output = []byte(strings.Join(conv.Rules, "\n"))
	default:
		return fmt.Errorf("unknown format %q", opt.Format)
	}

	// write to output file or setdout
	if out == "" {
		_, err = os.Stdout.Write(output)
		if err == nil {
			_, err = os.Stdout.WriteString("\n")
		}
	} else {
		err = ioutil.WriteFile(out, output, 0644)
	}

	return err
}
