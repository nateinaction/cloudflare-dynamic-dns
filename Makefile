
LINTERS := -E gci -E gofmt -E whitespace -E misspell -E gosec -E goconst

all: lint test build

lint:
	golangci-lint run $(LINTERS)

format:
	golangci-lint run --fix $(LINTERS)

test:
	go test ./... -cover

build:
	goreleaser release --snapshot --rm-dist

.PHONY: all lint format test build image