# coredns-operator

## Development
Pre-reqs:
`go`, `docker`, `operator-sdk`, `dep`, `kind`

```shell
make test

make kind-cluster
export KUBECONFIG="$(kind get kubeconfig-path --name="coredns-operator")"

make install
make run

make show
make install-cr
make show

# run in cluster
make docker-build
make kind-load
make deploy

make clean
```

## How this operator was constructed
Other than a few `dep` operations,
these directions were followed pretty closely:
https://github.com/operator-framework/getting-started

```shell
operator-sdk new coredns-operator
cd coredns-operator
operator-sdk add api --api-version=addons.k8s.io/v1alpha1 --kind=CoreDNS
operator-sdk generate k8s
operator-sdk add controller --api-version=addons.k8s.io/v1alpha1 --kind=CoreDNS
```
