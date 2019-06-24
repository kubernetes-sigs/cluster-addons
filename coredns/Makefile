
# Image URL to use all building/pushing image targets
IMG ?= ${USER}/coredns-operator:latest

all: test manager

# Run tests
test: generate fmt vet manifests
	go test ./pkg/... ./cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager sigs.k8s.io/addon-operators/coredns/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./cmd/manager/main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crds

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crds
	cat k8s/manager.yaml | sed -e "s@IMAGE_NAME@${IMG}@g" | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

push-image:
	mkdir -p .build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o .build/manager ./cmd/manager
	mkdir -p .build/bin
	cd .build/bin; wget -N https://storage.googleapis.com/kubernetes-release/release/v1.14.0/bin/linux/amd64/kubectl
	chmod +x .build/bin/kubectl
	tar cvf .build/layer.tar -C .build manager bin/kubectl
	tar rvf .build/layer.tar channels/
	crane append --base gcr.io/distroless/static:latest --new_layer .build/layer.tar --new_tag ${IMG} --output .build/merged.tar
	crane push .build/merged.tar ${IMG}
