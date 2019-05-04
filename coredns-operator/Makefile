
# Image URL to use all building/pushing image targets
TAG ?= latest
IMG ?= ${USER}/coredns-operator:${TAG}
NAMESPACE ?= "default"

all: test docker-build

# Run tests
test: fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	operator-sdk generate k8s

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	OPERATOR_NAME=coredns WATCH_NAMESPACE=${NAMESPACE} go run ./cmd/manager/main.go

# Start a kubernetes cluster that's missing CoreDNS
kind-cluster:
	kind create cluster --name coredns-operator
	export KUBECONFIG="$$(kind get kubeconfig-path --name="coredns-operator")"; \
	  kubectl -n kube-system delete deployment,configmap coredns; \
	  kubectl -n kube-system delete service kube-dns;

# Load the docker-build into the kind cluster so it can be run
kind-load:
	kind load docker-image --name coredns-operator ${IMG}

# Install CRDs into a cluster
install: manifests
	kubectl apply -f deploy/crds/addons_v1alpha1_coredns_crd.yaml

# Install CRs into a cluster
install-cr: manifests
	kubectl apply -f deploy/crds/addons_v1alpha1_coredns_cr.yaml

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f deploy/crds/addons_v1alpha1_coredns_crd.yaml
	kubectl apply -f deploy/clusterrole.yaml
	cat deploy/clusterrolebinding.yaml | sed -e "s@REPLACE_NAMESPACE@${NAMESPACE}@g" | kubectl apply -f -
	cat deploy/operator.yaml | sed -e "s@REPLACE_IMAGE@${IMG}@g" | kubectl -n ${NAMESPACE} apply -f -
	kubectl -n ${NAMESPACE} apply -f deploy/role.yaml
	kubectl -n ${NAMESPACE} apply -f deploy/role_binding.yaml
	kubectl -n ${NAMESPACE} apply -f deploy/service_account.yaml
	# Kick the deployment even when using the latest tag
	kubectl -n ${NAMESPACE} delete pods -l name=coredns-operator

# Show operator managed resources
show:
	kubectl -n ${NAMESPACE} get all,configmap,clusterrole,clusterrolebinding,serviceaccount,coredns -l k8s-app=kube-dns

# Build the docker image
docker-build: test
	operator-sdk build ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}

# Cleanup binaries, images, and cluster
clean: clean-bin clean-docker clean-kind
clean-bin:
	rm -rf ./build/_output
clean-docker:
	docker image rm ${IMG}
clean-kind:
	kind delete cluster --name coredns-operator
