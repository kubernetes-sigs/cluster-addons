# An RBAC Generator for Kubernetes Manifest

## Introduction

The rbac generator parses a kubernetes manifest and outputs the necessary rbac for an operator to apply and manage it.
The tool is inspired by the work being done on operators for cluster addons.

The tool works in two modes based on a boolean flag -   supervisory . 
When the flag is set to false, it parses a clusterrole in the manifest and adds the resources and verbs there to the
 generated cluster role. This is because an operator cannot create a clusterrole with permissions that it does not have.
 When the flag is set to true, it skips this and assumes that the clusterroles will be applied by something with
  higher permissions.

## Usage

1. Get the code
  
```shell script
go get https://github.com/kubernetes-sigs/cluster-addons
```
   
2. Build

```shell script
cd $GOPATH/kubernetes-sigs/cluster-addons/tools/rbac-gen
go install
```


- Run the go program

```shell script
rbac-gen --yaml <YOUR-MANIFEST-YAML>
```
Example:
```shell script
rbac-gen --yaml channels/packages/dashboard/manifest.yaml
```

There are optional flags like 
- `ns` for namespace: defaults to `kube-system`
- `out` - path to output file: output to stdout if not set
- `name` - name of the role to be generated

Have fun! :tada:
