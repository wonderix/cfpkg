
# Image URL to use all building/pushing image targets
OS := $(shell uname )
VERSION := $(shell git describe --tags --always --dirty)
REPOSITORY := wonderix/cfpkg
IMG ?= ${REPOSITORY}:${VERSION}
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd"
SHELL:=/bin/bash 

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GO_FILES:=$(shell git ls-files '*.go')

all: cfpkg

# Run tests
test: fmt vet 
	go test ./... -coverprofile cover.out

test-e2e: bin/cfpkg
	bin/cfpkg apply chart/worlds-simplest-service-broker
	bin/cfpkg test test/*_test.star

test-manual: bin/cfpkg
	bin/cfpkg test test/*_manual.star

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet 
	go run ./main.go


# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

cfpkg :: bin/cfpkg

VERSION_FLAGS := "-X github.com/wonderix/cfpkg/pkg/cfpkg.version=${VERSION}"

bin/cfpkg: $(GO_FILES) go.sum
	mkdir -p bin
	CGO_ENABLED=0 GOARCH=amd64 GO111MODULE=on go build -ldflags ${VERSION_FLAGS} -o bin/cfpkg . 

bin/linux/cfpkg: $(GO_FILES) go.sum
	mkdir -p bin/linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags ${VERSION_FLAGS} -o bin/linux/cfpkg . 

binaries:
	mkdir -p bin
	cd bin; \
	for GOOS in linux darwin windows; do \
	  CGO_ENABLED=0 GOOS=$$GOOS GOARCH=amd64 GO111MODULE=on go build -ldflags ${VERSION_FLAGS} -o cfpkg ..; \
		tar czf cfpkg-binary-$$GOOS.tgz cfpkg; \
	done
