COMMIT_HASH=$(shell git rev-parse --short HEAD)

build:
	go build -o ./bin/bk -ldflags '-X github.com/yoskeoka/bookkeeping.CommitHash=$(COMMIT_HASH)' ./cmd/bk
