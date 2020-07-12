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
	go build -ldflags "-s -w -X main.version=$(VERSION)" -v ./cmd/$(NAME)/...

.PHONY: test
test: ## run tests
	go test -v ./...

coverage.out:
	go test -covermode=count -coverprofile=coverage.out ./pkg/$(NAME)/...

.PHONY: cover
cover: coverage.out ## run test coverage
	go tool cover -html=coverage.out

.PHONY: clean
clean: ## clean up
	go clean
	rm -f $(NAME)
	rm -f coverage.out
	rm -rf release

.PHONY: format
format: ## format code
	go fmt -x ./...

.PHONY: run
run: build ## run for debug
	@$(CURDIR)/$(NAME) -c seq,name,bool $(CURDIR)/pkg/taon/testdata/data.json
	@echo
	@$(CURDIR)/$(NAME) $(CURDIR)/pkg/taon/testdata/misc-array.json
	@echo
	@cat $(CURDIR)/pkg/taon/testdata/all_docs.json | jq .rows | $(CURDIR)/$(NAME) -c key,doc._id,doc._rev,doc.name,doc.rank
	@echo
	@$(CURDIR)/$(NAME) $(CURDIR)/pkg/taon/testdata/long-field.json
	@echo

.PHONY: deps
deps: ## install deps
	go get -t ./...

.PHONY: release
release: windows linux darwin ## build binaries for release

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o release/$(NAME)-$(VERSION)-$(os)-amd64 && gzip -9 release/$(NAME)-$(VERSION)-$(os)-amd64
