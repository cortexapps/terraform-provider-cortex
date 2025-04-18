HOSTNAME=github.com
NAMESPACE=cortexapps
NAME=cortex
BINARY=terraform-provider-${NAME}
VERSION=0.4.3-dev

GOOS?=$(shell go env | grep GOOS | cut -d '=' -f2 | tr -d "'")
GOARCH?=$(shell go env | grep GOARCH | cut -d '=' -f2 | tr -d "'")
OS_ARCH?=$(GOOS)_$(GOARCH)

TF_LOG ?= WARN
include .env
export

default: test

# Run tests
.PHONY: test

build:
	go build -o ./bin/${BINARY} -ldflags="-X 'main.version=${VERSION}'"

release:
	export GOOS=linux GOARCH=amd64 && go build -o ./bin/${BINARY}_${VERSION}_$${GOOS}_$${GOARCH} -ldflags="-X 'main.version=${VERSION}'" && cd bin && zip ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH}.zip ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH} && rm ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH}
	export GOOS=darwin GOARCH=amd64 && go build -o ./bin/${BINARY}_${VERSION}_$${GOOS}_$${GOARCH} -ldflags="-X 'main.version=${VERSION}'" && cd bin && zip ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH}.zip ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH} && rm ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH}
	export GOOS=darwin GOARCH=arm64 && go build -o ./bin/${BINARY}_${VERSION}_$${GOOS}_$${GOARCH} -ldflags="-X 'main.version=${VERSION}'" && cd bin && zip ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH}.zip ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH} && rm ${BINARY}_${VERSION}_$${GOOS}_$${GOARCH}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ./bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

lint:
	golangci-lint run

docs:
	go generate ./...

format:
	go fmt ./...

# unit tests
test:
	go clean -testcache
	go test -v -cover ./... $(TESTARGS)

# acceptance tests
testacc:
	go clean -testcache
	TF_LOG=$(TF_LOG) TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m
