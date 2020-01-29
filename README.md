# Addon Operators

Repo for [cluster addon operators](https://github.com/kubernetes/enhancements/blob/master/keps/sig-cluster-lifecycle/addons/0035-20190128-addons-via-operators.md). More discussion can be seen in the [original KEP PR](https://github.com/kubernetes/enhancements/pull/746).

## Frequently asked questions

> What is this?

Cluster Addons is a sub-project of [SIG-Cluster-Lifecycle](https://github.com/kubernetes/community/tree/master/sig-cluster-lifecycle). Addon management has been a problem of cluster tooling for a long time.

This sub-project wants to figure out the scope of cluster addons. We want this to be manageable and not try to solve every single deployment problem.

> What is this not?

This sub-project is not interested in maintaining all cluster addons. Here we want to create some design patterns, some libraries, some supporting tooling, so everybody can easily create their own operators.

Not everything will need a cluster addon. Not everyone will want to use an operator.

> What is a cluster addon?

The lifecycle of a cluster addon is managed alongside the lifecycle of the cluster. Typically it has to be upgraded/downgraded when you move to a newer Kubernetes version. We want to use operators for this: a CRD describes the addon, and then the code which installs whatever the addon does, controlled by the CRD.

> What's your current agenda and timeline?

For now we want to create an actual [cluster addon we deemed as straight-forward](https://github.com/kubernetes-sigs/cluster-addons/issues/3), so we have actual code to look at, can learn from this experience (and what others have done in the past) and then take on some of the bigger philosophical questions.

> Who does this?

Cluster addons is a community project. If you're interested in building this, please get in touch. We're all ears!

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

Check out up to date information about where discussions and meetings happen on
the [community page of SIG Cluster Lifecycle](https://github.com/kubernetes/community/tree/master/sig-cluster-lifecycle).

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).
