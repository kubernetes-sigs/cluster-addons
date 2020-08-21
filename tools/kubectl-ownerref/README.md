# Kubectl Ownerrefs

## Introduction
`kubectl ownerref` is a simple kubectl plugin for getting all the resources in a cluster that another resource has ownerrefs on. It gives the basic information at a quick glance. It is used in cluster-addons to get the resources that a custom resource has ownerrefs on.

  ## Usage

  1. Get the code

  ```shell script
  go get https://github.com/kubernetes-sigs/cluster-addons
  ```

  2. Build

  ```shell script
  cd $GOPATH/kubernetes-sigs/cluster-addons/tools/kubectl-ownerref
  go install
  ```

  3. Run the go program

```shell script
kubectl ownerref <KIND> <NAME>
```

Example:
```shell script
kubectl ownerref Dashboard dashboard_sample
```

Have fun! :tada:
