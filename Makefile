export GO111MODULE=on
GO_SRC=$(shell find . -path ./.build -prune -false -o -name \*.go)

.PHONY: all
all: lint test

test: $(GO_SRC)
	go test -v -race -cover -coverpkg ./... -coverprofile=coverage.txt -covermode=atomic ./...

lint: ./.golangcilint.yaml
	./bin/golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.46.2 
	./bin/golangci-lint --config ./.golangcilint.yaml run ./...

.PHONY: check-licenses
check-licenses:
	go install github.com/google/go-licenses@latest
	@for tag in "$(EXTENSIONS),containers_image_openpgp" "$(EXTENSIONS),containers_image_openpgp"; do \
		echo Evaluating tag: $$tag;\
		for mod in $$(go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all); do \
			while [ x$$mod != x ]; do \
				echo -n "Checking $$mod ... "; \
				result=$$(GOFLAGS="-tags=$${tag}" go-licenses check $$mod 2>&1); \
				if [ $$? -eq 0 ]; then \
					echo OK; \
					break; \
				fi; \
				echo "$${result}" | grep -q "Forbidden"; \
				if [ $$? -eq 0 ]; then \
					echo FAIL; \
					exit 1; \
				fi; \
				echo "$${result}" | egrep -q "missing go.sum entry|no required module provides package|build constraints exclude all|updates to go.mod needed"; \
				if [ $$? -eq 0 ]; then \
					echo UNKNOWN; \
					break; \
				fi; \
			done; \
		 done; \
	 done
