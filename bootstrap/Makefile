# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DEV_TAG?=latest
DEV_REPO?=${USER}/cluster-addons-bootstrap
DEV_IMG?=$(DEV_REPO):$(DEV_TAG)
IMG?=$(DEV_IMG)
ARCH?=amd64

.PHONY: build

all: build

build:
	CGO_ENABLED=0 GOARCH=$(ARCH) GO111MODULE=on go build -o cluster-addons-bootstrap main.go

docker-image:
  # We must pre-pull images https://github.com/moby/buildkit/issues/1271
  # This also makes image building more understandable by avoiding a stale cache
	docker pull docker.io/library/ubuntu:20.04
	docker pull docker.io/library/golang:1.15

  # BuildKit is required for auto-population of `ARG TARGETARCH` during the build
  # BuildKit requires Docker >= 18.09
	DOCKER_BUILDKIT=1 docker build --platform=linux/$(ARCH) . -t $(IMG)

docker-push: docker-image
	docker push $(IMG)

fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test ./... -coverprofile cover.out

clean:
	docker rmi -f $(IMG)
