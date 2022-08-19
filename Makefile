.DEFAULT_GOAL := test
NAME := taon
SRC := $(wildcard ./cmd/$(NAME)/*.go ./pkg/$(NAME)/*.go)

.PHONY: build
build: $(NAME)

$(NAME): $(SRC)
	go build -o $(NAME) ./cmd/$(NAME)/...

.PHONY: test
test:
	go test -v ./pkg/taon/...

coverage.out:
	go test -covermode=count -coverprofile=coverage.out ./pkg/$(NAME)/...

.PHONY: cover
cover: coverage.out
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	go clean
	rm -f $(NAME)
	rm -f coverage.out
	rm -rf dist

.PHONY: run
run: $(NAME)
	./$(NAME) -c seq,name,bool $(CURDIR)/pkg/taon/testdata/data.json
	@echo
	./$(NAME) --columns seq,name,bool --markdown $(CURDIR)/pkg/taon/testdata/data.json
	@echo
	./$(NAME) $(CURDIR)/pkg/taon/testdata/misc-array.json
	@echo
	cat $(CURDIR)/pkg/taon/testdata/misc-array.json | ./$(NAME) -
	@echo
	cat $(CURDIR)/pkg/taon/testdata/all_docs.json | jq .rows | ./$(NAME) -c key,doc._id,doc._rev,doc.name,doc.rank
	@echo
	./$(NAME) $(CURDIR)/pkg/taon/testdata/long-field.json
	@echo

.PHONY: release
release:
	goreleaser check
	goreleaser build --snapshot --rm-dist
