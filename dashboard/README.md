# This is a simple guide for building an operator with kubebuilder version 2

## 1. Create a folder
```bash
mkdir dashboard
cd dashboard
```

## 2. Enable kubebuilder plugins and init
```bash
export KUBEBUILDER_ENABLE_PLUGINS=1
kubebuilder init --fetch-deps=false --domain=x-k8s.io --license=apache2
kubebuilder create api --pattern=addon --controller=true --example=false --group=addons --kind=Dashboard --make=false
 --namespaced=true --resource=true --version=v1alpha1
```

## 3. Run go mod vendor:

```bash
go mod vendor
```

## 3. Delete unnecessary test files
Delete the test suites that are more checking that kubebuilder is working:

```bash
find . -name "*_test.go" -delete
```

# Determine the a stable version and include its manifest:

Make a directory under `channels/packages/dashboard` with the version number. For the dashboard operator the stable
 version was `2.0.0`
```bash
cd  channels/packages/dashboard
mkdoir 2.0.0
```

Download the manifest yaml
```bash
cd 2.0.0
wget kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
mv recommended.yaml manifest.yaml
```

Define the stable channel:

```bash
cat > channels/stable <<EOF
manifests:
- version: 2.0.0
EOF
```

Running the operator locally:

```bash
make install
make run
```
Logs from the operator appear in the console.

You can now continue to build your operator by following the [kubebuilder-declarative-pattern](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#misc)


