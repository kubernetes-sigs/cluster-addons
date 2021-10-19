package main

import (
	"context"
	goflags "flag"
	"fmt"
	"io/ioutil"
	"os"

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

	return cmd
}

func RunGenerate(ctx context.Context, yamlFile string, out string, opt convert.BuildRoleOptions) error {
	//	read yaml file passed in from cmd
	bytes, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	// generate Group and Kind

	output, err := convert.ParseYAMLtoRole(string(bytes), opt)
	if err != nil {
		return err
	}

	if out == "" {
		fmt.Fprintf(os.Stdout, output)
	} else {
		err = ioutil.WriteFile(out, []byte(output), 0644)
	}

	return err
}
