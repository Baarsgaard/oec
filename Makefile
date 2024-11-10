.PHONY: all
all: fmt test build

.PHONY: build
build:
	go build -a -o oec

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
