.DEFAULT_GOAL := help

.PHONY: help
help: ## this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: print
print: ## print test data
	@cat $(PWD)/testdata/object.json | go run main.go
	@echo
	@cat $(PWD)/testdata/array.json | go run main.go

.PHONY: test
test: ## run tests
	@go test -v

.PHONY: fmt
fmt: ## format code
	@go fmt -x *.go
