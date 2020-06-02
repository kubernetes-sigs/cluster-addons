> Note: The CoreDNS operator is currently considered as `alpha`, and it is NOT recommended to be used in production.

# CoreDNS Operator

The CoreDNS Operator has been built to enable users to install the [CoreDNS](https://coredns.io) addon 
on their Kubernetes clusters

## Usage

The CoreDNS operator installs CoreDNS on the cluster helps to manage its resources
All the resources are installed via the use of `Kustomize`
This allows us to install the CoreDNS ConfigMap using the `configMapGenerator`, hashing the ConfigMap, 
which allows the CoreDNS deployment to undergo a proper and safe RollingUpdate


One of the main functionality of the operator is to be constantly watching the CoreDNS resources (deployment, ConfigMap, service etc.) and ensuring that it is in a functioning state. 
Any modification to the CoreDNS resources will result in the operator to reconcile and revert the changes

If there are any changes that is desired in CoreDNS, it can be done via the CoreDNS Custom Resource(CR)
The CR defines all the necessary specifications required by CoreDNS (example: CoreDNS Version, DNS Domain, Cluster IP and Corefile)

An example CR is as follows:

```yaml
apiVersion: addons.x-k8s.io/v1alpha1
kind: CoreDNS
metadata:
  name: coredns-operator
  namespace: kube-system
spec:
  version: 1.3.1
  dnsDomain: cluster.local
  dnsIP: 10.96.0.10
  corefile: |
    .:53 {
        errors
        health
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
```

The above CR will install CoreDNS version `1.3.1`, with DNS Domain `cluster.local`, Service IP `10.96.0.10` and the Corefile defined in the CR.

We can modify the specifications of CoreDNS by editing the Custom Resource.

For example, we can upgrade the CoreDNS version to `1.6.7` here by editing the `version` spec in the CR to `1.6.7`. 
This will enable the addon operator to install the manifests of CoreDNS associated with CoreDNS version `1.6.7`.

Another functionality that the operator provides while upgrading the CoreDNS version is the migration of the Corefile.
The operator will check if the existing Corefile is compatible with the new version of CoreDNS (In this case, from 1.3.1 -> 1.6.7) and will make changes accordingly.


> NOTE: While it is possible to downgrade the CoreDNS version, it is NOT recommended.

Currently, the operator can be used by running it locally outside the cluster, or we can also run the operator in-cluster

### Running the operator locally:

We can register the CoreDNS CRD and a CoreDNS object, and then try running
the controller locally.

1) We need to generate and register the CRDs:

```bash
$ make install
  /Users/srajan/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
  kustomize build config/crd | kubectl apply -f -
  customresourcedefinition.apiextensions.k8s.io/coredns.addons.x-k8s.io created
```

To verify that the CRD has registered successfully:

```bash
$ kubectl get crd coredns.addons.x-k8s.io
```

2) Create a CoreDNS CR:

```bash
$ kubectl apply -f config/samples/addons_v1alpha1_coredns.yaml 
coredns.addons.x-k8s.io/coredns-operator created

```

To verify that the CR has been created successfully:

```bash
$ kubectl get coredns -n kube-system
NAME               AGE
coredns-operator   3m54s
```

3) The controller can now be run using:

```bash
make run
```

We can see logs from the operator!

### Installing the operator in the cluster

To start, build the operator image:

```bash
make docker-build docker-push
```

Once the image has been built successfully, to build the CRD and start the operator:

```bash
make deploy
```

You can troubleshoot the operator by inspecting the controller:

```bash
$ kubectl -n coredns-operator-system get deploy
NAME                                  READY   UP-TO-DATE   AVAILABLE   AGE
coredns-operator-controller-manager   1/1     1            1           111s

# To check logs of the manager
$ kubectl -n coredns-operator-system logs <coredns-operator-controller-manager-pod-name> manager
```
