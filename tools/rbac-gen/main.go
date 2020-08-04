package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"sigs.k8s.io/cluster-addons/tools/rbac-gen/pkg/convert"
)

var (
	yamlFile    = flag.String("yaml", "manifest.yaml", "yaml file from which the rbac will be generated.")
	name        = flag.String("name", "generated-role", "name of role to be generated")
	saName      = flag.String("sa-name", "", "name of service account the role should be binded to")
	ns          = flag.String("ns", "kube-system", "namespace of the role to be generated")
	out         = flag.String("out", "", "name of output file")
	supervisory = flag.Bool("supervisory", false, "outputs role for operator in supervisory mode")
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()
	//	read yaml file passed in from cmd
	bytes, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		return err
	}

	// generate Group and Kind
	output, err := convert.ParseYAMLtoRole(string(bytes), *name, *ns, *saName, *supervisory)
	if err != nil {
		return err
	}

	if *out == "" {
		fmt.Fprintf(os.Stdout, output)
	} else {
		err = ioutil.WriteFile(*out, []byte(output), 0644)
	}

	return err
}
