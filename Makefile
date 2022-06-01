export GO111MODULE=on
GO_SRC=$(shell find . -path ./.build -prune -false -o -name \*.go)

.PHONY: all
all: lint test

test: $(GO_SRC)
	go test -v -race -cover -coverpkg ./... -coverprofile=coverage.txt -covermode=atomic ./...
