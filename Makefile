GO_SRC=$(shell find . -path ./.build -prune -false -o -name \*.go)
VERSION=$(shell git describe --tags || git rev-parse HEAD)
VERSION_FULL=$(if $(shell git status --porcelain --untracked-files=no),$(VERSION)-dirty,$(VERSION))

BUILD_TAGS = osusergo netgo static_build

build_zebra = go build -tags "$(BUILD_TAGS)" -buildmode=pie -ldflags "-X main.version=$(VERSION_FULL) -extldflags '-static'" -o zebra ./cmd/client
build_zebra_server = go build -tags "$(BUILD_TAGS)" -buildmode=pie -ldflags "-X main.version=$(VERSION_FULL) -extldflags '-static'" -o zebra-server ./cmd/server

zebra: $(GO_SRC) go.mod go.sum
	$(call build_zebra)

zebra-server: $(GO_SRC) go.mod go.sum
	$(call build_zebra_server)

bin: zebra zebra-server

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
	go fmt ./... && ([ -z $(CI) ] || git diff --exit-code)

test: $(GO_SRC)
	go test -v -race -cover -coverpkg ./... -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: check
check: test lint

.PHONY: mod
mod:
	go get -u
	go mod tidy

.PHONY: clean
clean:
	-rm zebra
	-rm zebra-server
	-rm *.crt
	-rm *.key
	-rm coverage.txt
	-rm -rf ./bin
