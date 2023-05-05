HOSTNAME=github.com
NAMESPACE=bigcommerce
NAME=cortex
BINARY=terraform-provider-${NAME}
VERSION=0.0.1-dev

GOOS?=$(shell go tool dist env | grep GOOS | grep -o '".*"' | sed 's/"//g')
GOARCH?=$(shell go tool dist env | grep GOARCH | grep -o '".*"' | sed 's/"//g')
OS_ARCH?=$(GOOS)_$(GOARCH)

TF_LOG ?= "WARN"
CORTEX_API_TOKEN ?= "set-me-in-env"

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

# unit tests
test:
	go clean -testcache
	go test -v -cover ./...

# acceptance tests
testacc:
	go clean -testcache
	CORTEX_API_TOKEN=$(CORTEX_API_TOKEN) TF_LOG=$(TF_LOG) TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
