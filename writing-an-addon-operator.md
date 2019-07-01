# Writing an Addon Operator

## What is it

The [Addons via Operators KEP](kep) discusses how operators can be used for managing cluster addons. The Cluster Addons sub-project has been discussing various approaches since the KEP was merged. Below we will discuss a simple example, which was the first addon operator that was reviewed and "blessed" by the team. It should be straight-forward and give you an idea of how to proceed writing your own addon operator.

## An example

Bringing up CoreDNS in a Kubernetes cluster was identified as a task that was clear and simple enough but still help us understand the general problem space and ask the right questions.

You can review the code of the [CoreDNS addon operator](coredns-op) in this repository. We will discuss the most interesting parts of it here. [Its README](coredns-readme) describes how the code scaffolding for the operator was set up, using `kubebuilder` and the [kubebuilder-declarative-pattern](kdp).

### XXX: Breakdown structure of operator

```sh
.
├── channels
│   ├── packages
│   │   └── coredns
│   │       └── 1.3.1
│   │           ├── coredns.yaml.base
│   │           └── manifest.yaml
│   └── stable
├── cmd
│   └── manager
│       └── main.go
├── config
│   ├── crds
│   │   └── addons_v1alpha1_coredns.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_image_patch.yaml
│   │   └── manager_prometheus_metrics_patch.yaml
│   ├── manager
│   │   └── manager.yaml
│   ├── rbac
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── manager_role_binding.yaml
│   │   └── manager_role.yaml
│   └── samples
│       └── addons_v1alpha1_coredns.yaml
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
├── k8s
│   └── manager.yaml
├── Makefile
├── pkg
│   ├── apis
│   │   ├── addons
│   │   │   ├── group.go
│   │   │   └── v1alpha1
│   │   │       ├── coredns_types.go
│   │   │       ├── doc.go
│   │   │       ├── register.go
│   │   │       └── zz_generated.deepcopy.go
│   │   ├── addtoscheme_addons_v1alpha1.go
│   │   └── apis.go
│   ├── controller
│   │   ├── add_coredns.go
│   │   ├── controller.go
│   │   └── coredns
│   │       └── coredns_controller.go
│   └── webhook
│       └── webhook.go
├── PROJECT
├── README.md
└── tools.go
```

### XXX: Lifecycle of operator

### XXX: Discussion of things strictly related to the coredns operator

### XXX: Running the operator

## Get started

XXX: Steal instructions from [coredns-readme].

## Talk to us

If you are interested in this, want to explore addon operators, want to discuss or get stuck somewhere, please talk to us:

- [#cluster-addons Slack](https://kubernetes.slack.com/messages/cluster-addons)
- [SIG Cluster Lifecycle group](https://groups.google.com/forum/#!forum/kubernetes-sig-cluster-lifecycle)

We appreciate your questions and feedback.

[kep]: https://github.com/kubernetes/enhancements/blob/master/keps/sig-cluster-lifecycle/addons/0035-20190128-addons-via-operators.md
[coredns-op]: https://github.com/kubernetes-sigs/addon-operators/tree/master/coredns
[coredns-readme]: https://github.com/kubernetes-sigs/addon-operators/blob/master/coredns/README.md
[kdp]: https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern
