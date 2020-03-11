# Example operator for CoreDNS

Broadly based on [kubebuilder-declarative-pattern walkthrough](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/blob/master/docs/addon/walkthrough/README.md)

A few differences so we can use go modules and [crane](https://github.com/google/go-containerregistry/blob/master/cmd/crane/doc/crane.md) - neither of which are required, just personal preference.

Created with kubebuilder:

```bash
export KUBEBUILDER_ENABLE_PLUGINS=1
kubebuilder init --fetch-deps=false --domain=x-k8s.io --license=apache2

kubebuilder create api --pattern=addon --controller=true --example=false --group=addons --kind=CoreDNS --make=false --namespaced=true --resource=true --version=v1alpha1

```

Run go mod vendor:

```bash
go mod vendor
```

Delete the test suites that are more checking that kubebuilder is working:

```bash
find . -name "*_test.go" -delete
```

Commit

```bash
git add .
git reset HEAD vendor
git commit -m "Initial CoreDNS scaffolding"
```



Create the manifests (we bake them into the addon-operator by default):

```bash
mkdir -p channels/packages/coredns/1.3.1/
pushd channels/packages/coredns/1.3.1/
wget https://raw.githubusercontent.com/kubernetes/kubernetes/9b437f95207c04bf2f25ef3110fac9b356d1fa91/cluster/addons/dns/coredns/coredns.yaml.base
cat coredns.yaml.base > manifest.yaml
popd
```

Define the stable channel:

```bash

cat > channels/stable <<EOF
manifests:
- version: 1.3.1
EOF

```

Running the operator locally:

```bash
make install
make run
```
We can see logs from the operator!


Generally follow the [main instructions](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/blob/master/docs/addon/walkthrough/README.md) at this point:

* [enable the declarative pattern library in your types](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#adding-the-framework-into-our-types) and
* [enable to declarative pattern in your controller](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#using-the-framework-in-the-controller)
* finally add the [call to addon.Init](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#misc)

Note that we intend to build these three steps into kubebuilder!

Then follow the instructions for deploying onto kubernetes.
