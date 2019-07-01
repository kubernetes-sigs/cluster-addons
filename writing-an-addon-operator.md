# Writing an Addon Operator

## What is it

The [Addons via Operators KEP](kep) discusses how operators can be used for managing cluster addons. The Cluster Addons sub-project has been discussing various approaches since the KEP was merged. Below we will discuss a simple example, which was the first addon operator that was reviewed and "blessed" by the team. It should be straight-forward and give you an idea of how to proceed writing your own addon operator.

## An example

Bringing up CoreDNS in a Kubernetes cluster was identified as a task that was clear and simple enough but still help us understand the general problem space and ask the right questions.

You can review the code of the [CoreDNS addon operator](coredns-op) in this repository. We will discuss the most interesting parts of it here. [Its README](coredns-readme) describes how the code scaffolding for the operator was set up, using `kubebuilder` and the [kubebuilder-declarative-pattern](kdp).

## Get started

## Talk to us

[kep]: https://github.com/kubernetes/enhancements/blob/master/keps/sig-cluster-lifecycle/addons/0035-20190128-addons-via-operators.md
[coredns-op]: https://github.com/kubernetes-sigs/addon-operators/tree/master/coredns
[coredns-readme]: https://github.com/kubernetes-sigs/addon-operators/blob/master/coredns/README.md
[kdp]: https://github.com/kubernetes-sigs/kubebuilder-declarative-pattern
