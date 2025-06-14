.DEFAULT_GOAL := test
NAME := taon
SRC := $(wildcard *.go ./pkg/$(NAME)/*.go)

.PHONY: build
build: $(NAME)

$(NAME): $(SRC)
	go build -o $(NAME) main.go

.PHONY: test
test:
	COLUMNS=158 go test -v ./pkg/taon/...

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
run: export COLUMNS=80
run: $(NAME)
	cat $(CURDIR)/pkg/taon/testdata/all_docs.json | jq .rows | ./$(NAME) -c key,doc._id,doc._rev,doc.name,doc.rank
	@echo
	cat $(CURDIR)/pkg/taon/testdata/all_docs.json | jq .rows[0] | ./$(NAME)

.PHONY: fixtures
fixtures: fixture_names=array data data_deep data_object long_field object
fixtures: $(NAME)
	$(foreach n,$(fixture_names),./$(NAME) ./pkg/taon/testdata/$(n).json > ./pkg/taon/testdata/$(n).txt;)
	$(foreach n,$(fixture_names),./$(NAME) --markdown ./pkg/taon/testdata/$(n).json > ./pkg/taon/testdata/$(n).md;)
	./$(NAME) -c seq,name,word ./pkg/taon/testdata/data.json > ./pkg/taon/testdata/data_columns.txt
	./$(NAME) -c key,value.rev,doc.name ./pkg/taon/testdata/data_deep.json > ./pkg/taon/testdata/data_deep_columns.txt
	./$(NAME) -c key,value.rev,doc.name ./pkg/taon/testdata/data_object.json > ./pkg/taon/testdata/data_object_columns.txt

.PHONY: release
release:
	goreleaser check
	goreleaser build --snapshot --clean
