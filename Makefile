
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

GO_FILES:=$(shell git ls-files '*.go') go.mod go.sum

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

VERSION_FLAGS := "-X github.com/wonderix/shalm/pkg/shalm.version=${VERSION}"

bin/cfpkg: $(GO_FILES)
	mkdir -p bin
	CGO_ENABLED=0 GOARCH=amd64 GO111MODULE=on go build -ldflags ${VERSION_FLAGS} -o bin/cfpkg . 

docker-context/cfpkg:  bin/linux/cfpkg
	mkdir -p docker-context/
	cp bin/linux/cfpkg docker-context/cfpkg

define BUILD_BIN
bin/$(1)/cfpkg: $(GO_FILES)  go.sum
	mkdir -p bin/$(1)
	CGO_ENABLED=0 GOOS=$(1) GOARCH=amd64 GO111MODULE=on go build -ldflags ${VERSION_FLAGS} -o bin/$(1)/cfpkg .
bin/cfpkg-binary-$(1).tgz: bin/$(1)/cfpkg
	cd bin/$(1) &&  tar czf ../cfpkg-binary-$(1).tgz cfpkg
endef

$(foreach i,linux darwin windows,$(eval $(call BUILD_BIN,$(i))))

binaries: $(foreach i,linux darwin windows,bin/cfpkg-binary-$(i).tgz)

formula: homebrew-tap/cfpkg.rb

homebrew-tap/cfpkg.rb: bin/cfpkg-binary-darwin.tgz bin/cfpkg-binary-linux.tgz
	@mkdir -p homebrew-tap
	@sed  \
	-e "s/{{sha256-darwin}}/$$(shasum -b -a 256 bin/cfpkg-binary-darwin.tgz  | awk '{print $$1}')/g" \
	-e "s/{{sha256-linux}}/$$(shasum -b -a 256 bin/cfpkg-binary-linux.tgz  | awk '{print $$1}')/g" \
	-e "s/{{version}}/$(VERSION)/g" homebrew-formula.rb \
	> homebrew-tap/cfpkg.rb
