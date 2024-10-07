.PHONY: all
all: fmt update test build

.PHONY: build
build:
	go build -a -o oec ./main

.PHONY: update
update:
	go get -u ./...
	go mod tidy
	go mod vendor

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: docker
docker:
	docker build -t oec:latest .

# Useful if more build dependencies are added later
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

export KO_DOCKER_REPO ?= ko.local/opsgenie/oec
.PHONY: ko
ko:
ifneq (,$(shell which ko))
	@go install github.com/google/ko@latest
endif
	$(GOBIN)/ko build --sbom=none --bare ./main
