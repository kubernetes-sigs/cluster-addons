# This is a simple guide for building an operator with kubebuilder version 2

Replace `dashboard` with name of the addon that you want to create.

## 1. Create a folder
```bash
mkdir dashboard
cd dashboard
```

## 2. Enable kubebuilder plugins and init
```bash
export KUBEBUILDER_ENABLE_PLUGINS=1
kubebuilder init --fetch-deps=false --domain=x-k8s.io --license=apache2
kubebuilder create api --pattern=addon --controller=true --example=false \
   --group=addons --kind=Dashboard --make=false \
   --namespaced=true --resource=true --version=v1alpha1
```

## 3. Run go mod vendor:

```bash
go mod vendor
```

## 3. Delete unnecessary test files
Delete the test suites that are checking that kubebuilder is working:

```bash
find . -name "*_test.go" -delete
```

# Determine the stable version and include its manifest:

Make a directory under `channels/packages/dashboard` for a stable version of the dashboard -- then download the manifest for it:
```bash
mkdir -p channels/packages/dashboard/2.0.0
wget -O channels/packages/dashboard/2.0.0/manifest.yaml https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0/aio/deploy/recommended.yaml
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

Now in another terminal, create an instance of the Dashboard custom resource:
```bash
kubectl create ns kubernetes-dashboard
kubectl -n kubernetes-dashboard apply -f config/samples/addons_v1alpha1_dashboard.yaml
```
You should see the operator respond and apply the resources from your package.

You can now continue to build your operator by following the [kubebuilder-declarative-pattern](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#misc)
