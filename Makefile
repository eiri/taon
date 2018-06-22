.DEFAULT_GOAL := all

NAME=taon

.PHONY: help
help: ## this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: all
all: deps test build ## test and build

.PHONY: build
build: ## build the binary
	go build -o $(NAME) -v

.PHONY: test
test: ## run tests
	go test -v ./...

.PHONY: clean
clean: ## clean up
	go clean
	rm -f $(NAME)

.PHONY: format
format: ## format code
	go fmt -x *.go

.PHONY: run
run: ## run for debug
	@cat $(PWD)/testdata/object.json | go run main.go
	@echo
	@cat $(PWD)/testdata/array.json | go run main.go

.PHONY: deps
deps: ## install deps
	go get -t ./...
