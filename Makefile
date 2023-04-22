BIN = bin
GO_FILES = $(shell find . -path '*/.*' -o -path '*/testdata/*' -prune \
	   -o '(' -type f -a -name '*.go' ')' -print)
REQUIREDFIELD = $(BIN)/requiredfield

REVIVE = $(BIN)/revive
STATICCHECK = $(BIN)/staticcheck

TOOLS = $(REVIVE) $(STATICCHECK)

PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
export GOBIN = $(PROJECT_ROOT)/$(BIN)

.PHONY: all
all: build lint test

.PHONY: build
build: $(REQUIREDFIELD)

$(REQUIREDFIELD): $(GO_FILES)
	go install go.abhg.dev/requiredfield/cmd/requiredfield

.PHONY: tools
tools: $(TOOLS)

.PHONY: test
test: $(GO_FILES)
	go test -v -race ./...

.PHONY: cover
cover: $(GO_FILES)
	go test -v -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint: gofmt revive staticcheck gomodtidy requiredfield

.PHONY: gofmt
gofmt:
	$(eval FMT_LOG := $(shell mktemp -t gofmt.XXXXX))
	@gofmt -e -s -l $(GO_FILES) > $(FMT_LOG) || true
	@[ ! -s "$(FMT_LOG)" ] || \
		(echo "gofmt failed. Please reformat the following files:" | \
		cat - $(FMT_LOG) && false)

.PHONY: requiredfield
requiredfield: $(REQUIREDFIELD)
	go vet -vettool=$(REQUIREDFIELD) ./...

.PHONY: revive
revive: $(REVIVE)
	$(REVIVE) -set_exit_status ./...

$(REVIVE): tools/go.mod
	go install -C tools github.com/mgechev/revive

.PHONY: staticcheck
staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

$(STATICCHECK): tools/go.mod
	go install -C tools honnef.co/go/tools/cmd/staticcheck

.PHONY: gomodtidy
gomodtidy: go.mod go.sum tools/go.mod tools/go.sum
	go mod tidy
	go mod tidy -C tools
	@if ! git diff --quiet $^; then \
		echo "go mod tidy changed files:" && \
		git status --porcelain $^ && \
		false; \
	fi
