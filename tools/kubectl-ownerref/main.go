package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/utils/pointer"

	// Load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	klog.InitFlags(nil)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()

	warningHandler := &warningHandler{
		fallback: rest.NewWarningWriter(os.Stderr, rest.WarningWriterOptions{Deduplicate: true}),
	}
	kubeConfigFlags.WrapConfigFn = func(c *rest.Config) *rest.Config {
		c.WarningHandler = warningHandler
		c.QPS = 100
		return c
	}

	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)

	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)

	ioStreams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}

	rootCommand := NewCmdGetByOwner("kubectl", f, ioStreams)
	rootCommand.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	matchVersionKubeConfigFlags.AddFlags(rootCommand.PersistentFlags())
	kubeConfigFlags.AddFlags(rootCommand.PersistentFlags())

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}

type GetByOwnerOptions struct {
	PrintFlags *get.PrintFlags
	ToPrinter  func(*meta.RESTMapping, *bool, bool, bool) (printers.ResourcePrinterFunc, error)

	Namespace         string
	ExplicitNamespace bool

	AllNamespaces bool

	OwnerKindArg string
	OwnerKind    schema.GroupKind
	OwnerName    string

	genericclioptions.IOStreams
}

type warningHandler struct {
	fallback rest.WarningHandler
}

func (h *warningHandler) HandleWarningHeader(code int, agent string, text string) {
	if code == 299 {
		// Ignore deprecated kind messages, the user isn't normally specifying the kinds
		return
	}
	h.fallback.HandleWarningHeader(code, agent, text)
}

// NewGetByOwnerOptions returns a GetByOwnerOptions with default options.
func NewGetByOwnerOptions(parent string, streams genericclioptions.IOStreams) *GetByOwnerOptions {
	return &GetByOwnerOptions{
		PrintFlags: get.NewGetPrintFlags(),

		IOStreams: streams,
	}
}

func NewCmdGetByOwner(parent string, f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewGetByOwnerOptions(parent, streams)

	cmd := &cobra.Command{
		Use: "kubectl ownerref <KIND> <NAME>",
		Example: `kubectl ownerref Dashboard dashboard_sample
kubectl ownerref Dashboard dashboard_sample -o yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd))
			cmdutil.CheckErr(o.Run(f, cmd))
		},
	}

	o.PrintFlags.AddFlags(cmd)

	cmd.Flags().BoolVarP(&o.AllNamespaces, "all-namespaces", "A", o.AllNamespaces, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	return cmd
}

// Complete takes the command arguments and factory and infers any remaining options.
func (o *GetByOwnerOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {

	if len(args) == 2 {
		o.OwnerKindArg = args[0]
		o.OwnerName = args[1]
	} else {
		return fmt.Errorf("syntax: kubectl ownerref <KIND> <NAME>")
	}

	var err error
	o.Namespace, o.ExplicitNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return err
	}
	if o.AllNamespaces {
		o.ExplicitNamespace = false
	}

	o.ToPrinter = func(mapping *meta.RESTMapping, outputObjects *bool, withNamespace bool, withKind bool) (printers.ResourcePrinterFunc, error) {
		// make a new copy of current flags / opts before mutating
		printFlags := o.PrintFlags.Copy()

		if mapping != nil {
			printFlags.SetKind(mapping.GroupVersionKind.GroupKind())
		}
		if withNamespace {
			printFlags.EnsureWithNamespace()
		}
		if withKind {
			printFlags.EnsureWithKind()
		}

		resourcePrinter, err := printFlags.ToPrinter()
		if err != nil {
			return nil, err
		}

		printer, err := printers.NewTypeSetter(scheme.Scheme).WrapToPrinter(resourcePrinter, nil)
		if err != nil {
			return nil, err
		}

		return printer.PrintObj, nil
	}

	return nil
}

// Validate checks the set of flags provided by the user.
func (o *GetByOwnerOptions) Validate(cmd *cobra.Command) error {
	return nil
}

// Run performs the get operation.
func (o *GetByOwnerOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	ctx := context.TODO()
	discoveryClientSet, err := f.ToDiscoveryClient()
	if err != nil {
		return err
	}

	clientset, err := f.DynamicClient()
	if err != nil {
		return err
	}

	apiResourceList, err := discoveryClientSet.ServerPreferredResources()
	if err != nil {
		return nil
	}

	// Find and normalize OwnerKind
	if o.OwnerKind.Kind == "" {
		var matches []schema.GroupKind
		for _, apiResource := range apiResourceList {
			for _, resource := range apiResource.APIResources {
				if strings.EqualFold(resource.Kind, o.OwnerKindArg) {
					gv, err := schema.ParseGroupVersion(apiResource.GroupVersion)
					if err != nil {
						return err
					}
					matches = append(matches, schema.GroupKind{Group: gv.Group, Kind: resource.Kind})
				}
			}
		}

		if len(matches) == 0 {
			return fmt.Errorf("kind %q not recognized", o.OwnerKindArg)
		}
		if len(matches) > 1 {
			return fmt.Errorf("kind %q is ambiguous", o.OwnerKindArg)
		}
		o.OwnerKind = matches[0]
	}

	// Walk all the kinds
	var errors []error
	for _, apiResource := range apiResourceList {
		for _, resource := range apiResource.APIResources {
			// Skip "fake" resources that don't support list
			hasList := false
			for _, verb := range resource.Verbs {
				if verb == "list" {
					hasList = true
				}
			}
			if !hasList {
				continue
			}

			gv, err := schema.ParseGroupVersion(apiResource.GroupVersion)
			if err != nil {
				return err
			}
			mapping := meta.RESTMapping{
				Resource: schema.GroupVersionResource{
					Group:    gv.Group,
					Version:  gv.Version,
					Resource: resource.Name,
				},
				GroupVersionKind: schema.GroupVersionKind{
					Group:   gv.Group,
					Version: gv.Version,
					Kind:    resource.Kind,
				},
			}
			if resource.Namespaced {
				mapping.Scope = meta.RESTScopeNamespace
			} else {
				mapping.Scope = meta.RESTScopeRoot
			}

			errs := o.PrintResource(ctx, clientset, &mapping)
			errors = append(errors, errs...)
		}
	}

	if len(errors) != 0 {
		for _, err := range errors {
			fmt.Fprintf(o.ErrOut, "error: %v\n", err)
		}
	}

	return nil
}

func (o *GetByOwnerOptions) PrintResource(ctx context.Context, c dynamic.Interface, mapping *meta.RESTMapping) []error {
	var errors []error
	var resource dynamic.ResourceInterface
	if mapping.Scope == meta.RESTScopeNamespace && !o.AllNamespaces {
		resource = c.Resource(mapping.Resource).Namespace(o.Namespace)
	} else {
		resource = c.Resource(mapping.Resource)
	}
	listOpts := metav1.ListOptions{}
	resourceList, err := resource.List(ctx, listOpts)
	if err != nil {
		errors = append(errors, fmt.Errorf("List(%v) failed: %w", mapping.Resource, err))
		return errors
	}

	return o.PrintList(mapping, resourceList)
}

// PrintList is a function for printing all kinds of resources
func (o *GetByOwnerOptions) PrintList(mapping *meta.RESTMapping, list *unstructured.UnstructuredList) []error {
	var errors []error

	if len(list.Items) == 0 {
		return errors
	}

	var ownerItems []unstructured.Unstructured
	for _, item := range list.Items {
		for _, ownerRef := range item.GetOwnerReferences() {
			if ownerRef.Kind != o.OwnerKind.Kind || ownerRef.Name != o.OwnerName {
				continue
			}
			gv, err := schema.ParseGroupVersion(ownerRef.APIVersion)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			if gv.Group == o.OwnerKind.Group {
				ownerItems = append(ownerItems, item)
			}
		}
	}

	if len(ownerItems) == 0 {
		return errors
	}

	printWithNamespace := !o.ExplicitNamespace
	printWithKind := true

	printer, err := o.ToPrinter(mapping, nil, printWithNamespace, printWithKind)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	for _, ownerItem := range ownerItems {
		if err := printer.PrintObj(&ownerItem, o.Out); err != nil {
			errors = append(errors, err)
		}
	}

	outputFormat := strings.ToLower(pointer.StringDeref(o.PrintFlags.OutputFormat, ""))
	switch outputFormat {
	case "":
		if _, err := o.Out.Write([]byte("\n")); err != nil {
			errors = append(errors, err)
		}

	case "yaml":
		if _, err := o.Out.Write([]byte("\n---\n")); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
