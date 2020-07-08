# Writing an Addon Operator

## What is it

The [Addons via Operators KEP](kep) details how operators can be used for managing cluster addons. Below we will present a simple example, It should be straight-forward and give you an idea of how to proceed writing your own addon operator.

## An example

Bringing up CoreDNS in a Kubernetes cluster had been identified as a task that was clear and simple enough but still help us understand the general problem space and ask the right questions.
You can review the code of the [CoreDNS addon operator](https://github.com/kubernetes-sigs/cluster-addons/tree/master/coredns) in this repository. 
Here we will take you through the most interesting parts. [Its README](https://github.com/kubernetes-sigs/cluster-addons/tree/master/coredns/README.md) describes how the code scaffolding for the operator was set up, using `kubebuilder` and the [kubebuilder-declarative-pattern](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern).

## Creating the operator

Broadly based on [kubebuilder-declarative-pattern walkthrough](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/blob/master/docs/addon/walkthrough/README.md)
A few differences so we can use go modules and [crane](https://github.com/google/go-containerregistry/blob/master/cmd/crane/doc/crane.md) - neither of which are required, just personal preference.

1. Created with kubebuilder:

```bash
export KUBEBUILDER_ENABLE_PLUGINS=1
kubebuilder init --fetch-deps=false --domain=x-k8s.io --license=apache2

kubebuilder create api --pattern=addon --controller=true --example=false --group=addons --kind=<my-addon> --make=false --namespaced=true --resource=true --version=v1alpha1

```

2. Run go mod vendor:

```bash
go mod vendor
```

3. Delete the test suites that are checking whether kubebuilder is working:

```bash
find . -name "*_test.go" -delete
```

4. Commit

```bash
git add .
git reset HEAD vendor
git commit -m "Initial addon scaffolding"
```

5. Create the manifests (we bake them into the addon-operator by default):

```bash
mkdir -p channels/packages/coredns/1.3.1/
pushd channels/packages/coredns/1.3.1/
wget https://raw.githubusercontent.com/kubernetes/kubernetes/9b437f95207c04bf2f25ef3110fac9b356d1fa91/cluster/addons/dns/coredns/coredns.yaml.base
cat coredns.yaml.base > manifest.yaml
popd
```

6. Define the stable channel:

```bash

cat > channels/stable <<EOF
manifests:
- version: 1.3.1
EOF

```

7. Generally follow the [main instructions](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/blob/master/docs/addon/walkthrough/README.md) at this point:

* [enable the declarative pattern library in your types](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#adding-the-framework-into-our-types) and
* [enable to declarative pattern in your controller](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#using-the-framework-in-the-controller)
* finally add the [call to addon.Init](https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern/tree/master/docs/addon/walkthrough#misc)

Note that we intend to build these three steps into kubebuilder!

Then follow the instructions for deploying onto kubernetes.

7. Running the operator locally:

```bash
make install
make run
```
We can see logs from the operator!

### Breakdown structure of operator

This is the structure of the operator after being created with Kubebuilder v2. The structure of any created operator will be similar to the one shown below
This example is of the [CoreDNS addon operator](https://github.com/kubernetes-sigs/cluster-addons/tree/master/coredns)

```sh
.
├── api
│   ├── v1alpha1
│   │   └── coredns
│   │       ├── coredns_types.go
│   │       ├── groupversion.go
│   │       └── zz_generated.deepcopy.go
├── channels
│   ├── packages
│   │   └── coredns
│   │       ├── 1.3.1
│   │       │   ├── clusterrole.yaml
│   │       │   ├── clusterrolebinding.yaml
│   │       │   ├── Corefile
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   ├── service.yaml
│   │       │   └── serviceaccount.yaml
│   │       ├── 1.6.7
│   │       │   ├── clusterrole.yaml
│   │       │   ├── clusterrolebinding.yaml
│   │       │   ├── Corefile
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   ├── service.yaml
│   │       │   └── serviceaccount.yaml
│   │       └── 1.6.9
│   │       │   ├── clusterrole.yaml
│   │       │   ├── clusterrolebinding.yaml
│   │       │   ├── Corefile
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   ├── service.yaml
│   │       │   └── serviceaccount.yaml
│   └── stable
├── config
│   ├── certmanager
│   │   └── certificate.yaml
│   │   └── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── crds
│   │   ├── bases
│   │   │   └── addons.x-k8s.io_coredns.yaml
│   │   ├── patches
│   │   │   └── cainjection_in_coredns.yaml
│   │   │   └── webhook_in_coredns.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_resource_patch.yaml
│   │   ├── manager_webhook_patch.yaml
│   │   └── webhookcainjection_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── coredns_editor_role.yaml
│   │   ├── coredns_viewer_role.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   ├── role.yaml
│   │   └── role_binding.yaml
│   ├── samples
│   │   └── addons_v1alpha1_coredns.yaml
│   └── webhook
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       └── manager.yaml
├── controllers
│   ├── tests
│   │   ├── patches-stable.in.yaml
│   │   ├── patches-stable.out.yaml
│   │   ├── simple-stable.in.yaml
│   │   └── simple-stable.out.yaml
│   ├── coredns_controller.go
│   ├── coredns_controller_test.go
│   └── util.go
├── hack
│   ├── smoketest.go
│   └── boilerplate.go.txt
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── Makefile
├── PROJECT
└── README.md
```

## Talk to us

If you are interested in this, want to explore addon operators, want to discuss or get stuck somewhere, you can contact us at:

- [#cluster-addons Slack](https://kubernetes.slack.com/messages/cluster-addons)
- [SIG Cluster Lifecycle group](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)
