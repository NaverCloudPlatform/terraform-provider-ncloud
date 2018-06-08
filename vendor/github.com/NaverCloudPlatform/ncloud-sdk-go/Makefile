GOFMT_FILES?=$$(find . -not -path "./vendor/*" -name "*.go")

ifeq ($(OS)),Windows_NT)
	SDK_ONLY_PKGS=$(shell go list ./... | findstr /v "\/vendor")
else
	SDK_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")
endif


all: deps test build

help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  build                   to go build the SDK"
	@echo "  test                    to run unit tests"

test:
	@echo "go(ginkgo) test SDK"
	@ginkgo -r

fmt:
	@gofmt -w -s $(GOFMT_FILES)

build:
	@echo "go build SDK"
	@go build ${SDK_ONLY_PKGS}

deps:
	@go get -u github.com/kardianos/govendor
	@govendor fetch +out

updatedeps:
	@govendor update +vendor
