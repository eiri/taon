.DEFAULT_GOAL := all

NAME := taon
VERSION := $(shell git describe --tags)
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: help
help: ## this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: deps test build ## test and build

.PHONY: build
build: ## build the binary
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(NAME) -v

.PHONY: test
test: ## run tests
	go test -v ./...

.PHONY: clean
clean: ## clean up
	go clean
	rm -f $(NAME)
	rm -rf release

.PHONY: format
format: ## format code
	go fmt -x *.go

.PHONY: run
run: ## run for debug
	@cat $(PWD)/testdata/object.json | go run main.go columns_value.go
	@echo
	@cat $(PWD)/testdata/array.json | go run main.go columns_value.go -c seq,number,name

.PHONY: deps
deps: ## install deps
	go get -t ./...

.PHONY: release
release: windows linux darwin ## build binaries for release

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o release/$(NAME)-$(VERSION)-$(os)-amd64 && gzip -9 release/$(NAME)-$(VERSION)-$(os)-amd64
