# Example operator for CoreDNS

Created with kubebuilder:

```bash
kubebuilder init --dep=false --domain=k8s.io --license apache2

kubebuilder create api --controller=true --example=false --group=addons --kind=CoreDNS --make=false --namespaced=true --resource=true --version=v1alpha1

```

Switched to go modules:

```bash
export GO111MODULE=on
go mod init sigs.k8s.io/addon-operators/coredns

# Insert our tools.go for extra dependencies
cp ../tools.go tools.go

go get -m k8s.io/client-go@v10.0.0
go get -m k8s.io/api@kubernetes-1.13.5
go get -m k8s.io/apimachinery@kubernetes-1.13.5
go get -m k8s.io/apiserver@kubernetes-1.13.5
go get -m k8s.io/apiextensions-apiserver@kubernetes-1.13.5

go mod vendor

rm Gopkg.toml
```

Delete the test suites that are more checking that kubebuilder is working:

```bash
find . -name "*_test.go" -delete
```

Commit

```bash
git add .
git reset HEAD vendor
```
