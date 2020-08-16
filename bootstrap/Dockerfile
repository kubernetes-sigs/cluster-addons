# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

# Download a supported `kubectl` release for the target arch
FROM --platform=$BUILDPLATFORM ubuntu:20.04 as kubectl
ARG TARGETARCH
RUN apt-get update
RUN apt-get install -y curl
RUN curl -fsSL https://dl.k8s.io/release/v1.17.4/bin/linux/${TARGETARCH}/kubectl > /usr/bin/kubectl
RUN ! ldd /usr/bin/kubectl # Assert that the downloaded kubectl is statically linked
RUN chmod a+rx /usr/bin/kubectl

# Build the bootstrap binary
FROM --platform=$BUILDPLATFORM golang:1.15 as builder

# Copy the Go Modules manifests
WORKDIR /go/src/bootstrap
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source files
COPY main.go main.go
COPY app/ app/

# Build
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} GO111MODULE=on go build -a -o cluster-addons-bootstrap main.go
RUN ! ldd cluster-addons-bootstrap # Assert that the compiled bin is statically linked

# Use distroless as minimal base image to package the bootstrap binary
# See: https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static
WORKDIR /
COPY --from=builder /go/src/bootstrap/cluster-addons-bootstrap .
COPY --from=kubectl /usr/bin/kubectl /usr/bin/kubectl

ENTRYPOINT ["/cluster-addons-bootstrap"]
