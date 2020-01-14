# addon-discovery

Identifies a cluster addon as a set of component Kubernetes resources selected by a label unique to that addon.

## Component Tracking

`addon-discovery` introduces a cluster-scoped `Addon` resource that surfaces:

- info about the addon
- the label selector used to gather its components
- a set of status enriched references to its components
- top-level conditions that summarize any abnormal state (Future Work)

Ex.

```yaml
apiVersion: discovery.addons.x-k8s.io/v1alpha1
kind: Addon
metadata:
  name: plumbus
status:
  components:
    matchLabels:
        discovery.addons.x-k8s.io/plumbus: ""
    refs:
    - kind: CustomResourceDefinition
      name: plumbai.how.dotheydoit.com
      apiVersion: apiextensions.k8s.io/v1beta1
      conditions:
      - lastTransitionTime: "2019-11-25T12:43:26Z"
        message: no conflicts found
        reason: NoConflicts
        status: "True"
        type: NamesAccepted
    - kind: Deployment
      name: plumbus-addon-controller
      apiVersion: apps/v1
      conditions:
      - lastTransitionTime: "2019-11-25T12:43:27Z"
        lastUpdateTime: "2019-11-25T12:43:39Z"
        message: ReplicaSet "plumbus-addon-controller-6999db5767" has successfully progressed.
        reason: NewReplicaSetAvailable
        status: "True"
        type: Progressing
   - kind: ClusterRoleBinding
     namespace: operators
     name: rb-9oacj
     apiVersion: rbac.authorization.k8s.io/v1
   # ...
```

In this case, any resource that bears the `discovery.addons.x-k8s.io/plumbus` label key will be tracked as a component of the `plumbus` `Addon`. In general, the component label key convention is `discovery.addons.x-k8s.io/<name>`, where `<name>` is the name of the desired addon. Additionally, there is no limit to the number of `Addons` that a particular resource is a component of.

## Automatic `Addon` Generation

When `addon-discovery` finds a resource bearing a component label for an `Addon` that doesn't exist yet, it automatically generates that `Addon`. This makes opt-in simple, only requiring users to apply labels.

## Metadata Extraction (Future Work)

Provide a way for an addon component to register a subset of data as interesting to the top-level addon.

The example below uses an annotation that indicates how to extract the GVK from CRDs that are part of the Addon:

```yaml
apiVersion: discovery.addons.x-k8s.io/v1alpha1
kind: Addon
metadata:
  name: plumbus
status:
  metadata:
    apis:
    - group: how.theydoit.com
      version: v2alpha1
      kind: Plumbus
      plural: plumbai
  components:
    # ...
---
kind: CRD
metadata:
    label:
        discovery.addons.x-k8s.io/plumbus: ""
    annotations:
        # jq object builder syntax?
        discovery.addons.x-k8s.io/metadata.apis: "[{group: spec.group, version: spec.versions[].name, kind: spec.names.kind, plural: spec.names.plural}]"
```

## Condition Probes (Future Work)

Add a resource that maps a query/response tuple to an output condition on its status. This resource can be included as a component of an `Addon` to drive its `status.conditions` field. Future iterations could support executing `Jobs` and parsing their results.

## Development

### Build and push the controller manager container image

```sh
$ IMG=quay.io/njhale/addon-controller:latest make build
# ...
$ IMG=quay.io/njhale/addon-controller:latest make docker-push
#...
```

__Note:__ _The `IMG` variable sets the output image reference; `quay.io/njhale/addon-controller:latest` should generally be changed to a container image registry that the developer has push access to._

### Deploy to the current kubectl context

```sh
$ IMG=quay.io/njhale/addon-controller@sha256:abcd make deploy
# ...
```

__Note:__ _The image referenced by `IMG` must be pullable by the target cluster._
