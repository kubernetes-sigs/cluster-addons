# This is the Kubernetes Dashboard Operator

This is an operator for managing the kubernetes dashboard addon.

## Available versions of kubernetes dashboard
```
Versions:
- 1.8.3
- 1.10.1
- 2.0.0
```

Stable version - `2.0.0`

## Usage
1. Clone the cluster-addons
```
go get https://github.com/kubernetes-sigs/cluster-addons
cd dashboard
```

2. Install Dashboard CRD
```
make install
```

3. Create Custom Resource. An example lives [here](https://raw.githubusercontent.com/kubernetes-sigs/cluster-addons/master/dashboard/config/samples/addons_v1alpha1_dashboard.yaml)

4. Run the operator locally
`make run`

or in-cluster
`make deploy`

5. View resources created in the `kubernetes-dashboard` namespace.
