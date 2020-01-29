SHELL := /bin/bash
PROJECT = sigs.k8s.io/cluster-addons/installer
APIS_DIR = ${PROJECT}/pkg/apis
CACHE_DIR = /tmp/go-cache

all: build

build: 
	$(MAKE) shell COMMAND="make binary"

shell:
	mkdir -p $(CACHE_DIR)/bin $(CACHE_DIR)/src $(CACHE_DIR)/cache bin/cache
	docker run -it \
		-v $(shell pwd):/go/src/sigs.k8s.io/cluster-addons/installer \
		-v $(CACHE_DIR)/bin:/go/bin \
		-v $(CACHE_DIR)/src:/go/src \
		-v $(CACHE_DIR)/cache:/.cache/go-build \
		-w /go/src/sigs.k8s.io/cluster-addons/installer \
		-u $(shell id -u):$(shell id -g) \
		-e GO111MODULE=on \
		golang:1.11 \
		$(COMMAND)

binary: autogen vendor
	go build -mod=vendor -o bin/installer sigs.k8s.io/cluster-addons/installer/cmd/installer

autogen: go/bin/deepcopy-gen go/bin/defaulter-gen go/bin/conversion-gen
	# Let the boilerplate be empty
	touch /tmp/boilerplate
	deepcopy-gen \
		--input-dirs ${APIS_DIR}/config,${APIS_DIR}/config/v1alpha1 \
		--bounding-dirs ${APIS_DIR} \
		-O zz_generated.deepcopy \
		-h hack/boilerplate.go.txt

	defaulter-gen \
		--input-dirs ${APIS_DIR}/config/v1alpha1 \
		-O zz_generated.defaults \
		-h hack/boilerplate.go.txt

	conversion-gen \
		--input-dirs ${APIS_DIR}/config,${APIS_DIR}/config/v1alpha1 \
		-O zz_generated.conversion \
		-h hack/boilerplate.go.txt

go/bin/%: vendor
	go install k8s.io/code-generator/cmd/$*

vendor:
	if [[ ! -f go.mod ]]; then go mod init; fi
	go mod tidy
	go mod vendor
	go mod verify

only-build:
	$(MAKE) shell COMMAND="go build -mod=vendor -o bin/installer sigs.k8s.io/cluster-addons/installer/cmd/installer"

clean:
	rm -rf bin vendor go.sum go.mod
	find . -type f | grep zz_generated | grep -v vendor | xargs -r rm
