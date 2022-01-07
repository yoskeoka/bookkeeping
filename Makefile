COMMIT_HASH=$(shell git rev-parse --short HEAD)
GOBIN=$(CURDIR)/bin

dev-tools:
	go install github.com/goreleaser/goreleaser@latest

build:
	go build -o ./bin/bk -ldflags '-X github.com/yoskeoka/bookkeeping.CommitHash=$(COMMIT_HASH)' ./cmd/bk
