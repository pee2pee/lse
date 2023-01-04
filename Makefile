HASGOCILINT := $(shell which golangci-lint 2> /dev/null)
ifdef HASGOCILINT
    GOLINT=golangci-lint
else
    GOLINT=bin/golangci-lint
endif

# Dependency versions
GOLANGCI_VERSION = 1.50.0

install:
	go install -v github.com/pee2pee/lse

build:
	go build -o ./bin/ ./.

test:
	go test -race ./...

.PHONY: fix
fix: ## Fix lint violations
	gofmt -s -w .
	goimports -w $$(find . -type f -name '*.go' -not -path "*/vendor/*")