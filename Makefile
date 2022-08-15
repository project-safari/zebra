GO_SRC=$(shell find . -path ./.build -prune -false -o -name \*.go)
VERSION=$(shell git describe --tags || git rev-parse HEAD)
VERSION_FULL=$(if $(shell git status --porcelain --untracked-files=no),$(VERSION)-dirty,$(VERSION))

BUILD_TAGS = osusergo netgo static_build

build_zebra = go build -tags "$(BUILD_TAGS)" -buildmode=pie -ldflags "-X main.version=$(VERSION_FULL) -extldflags '-static'" -o zebra ./cmd/client
build_zebra_server = go build -tags "$(BUILD_TAGS)" -buildmode=pie -ldflags "-X main.version=$(VERSION_FULL) -extldflags '-static'" -o zebra-server ./cmd/server
build_herd = go build -tags "$(BUILD_TAGS)" -buildmode=pie -ldflags "-X main.version=$(VERSION_FULL) -extldflags '-static'" -o herd ./cmd/herd

zebra: $(GO_SRC) go.mod go.sum
	$(call build_zebra)

zebra-server: $(GO_SRC) go.mod go.sum
	$(call build_zebra_server)

herd: $(GO_SRC) go.mod go.sum
	$(call build_herd)

bin: zebra zebra-server herd

lint: ./.golangcilint.yaml
	./bin/golangci-lint --version || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.46.2 
	./bin/golangci-lint --config ./.golangcilint.yaml run ./...

.PHONY: simulator-setup
simulator-setup: bin
	rm -rf ./simulator/simulator-store && ./herd --store ./simulator/simulator-store
	rm -f ./simulator/zebra-simulator.json
	rm -f ./simulator/admin.yaml
	./zebra -c ./simulator/admin.yaml config init https://127.0.0.1:6666
	./zebra -c ./simulator/admin.yaml config email admin@zebra.project-safari.io
	./zebra -c ./simulator/admin.yaml config ca-cert ./simulator/zebra-ca.crt
	./zebra-server -c ./simulator/zebra-simulator.json init --auth-key "AvadaKedavra" --user="./simulator/admin.yaml" --password "Riddikulus" --cert "./simulator/zebra-server.crt" --key "./simulator/zebra-server.key" -a "tcp://127.0.0.1:6666" --store="./simulator/simulator-store"
	sed -i 's/ravi/admin/g' ./simulator/admin.yaml

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

.PHONY: certs
certs:
	./simulator/gen_certs.sh "./simulator" "zebra.inisieme.local" "1.1.1.1"
	cd ./simulator && go test

simulator: bin certs simulator-setup
	./zebra-server --config ./simulator/zebra-simulator.json

test: $(GO_SRC) certs simulator-setup
	go test -v -race -cover -coverpkg ./... -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: check
check: test lint

.PHONY: mod
mod:
	go get -u
	go mod tidy

.PHONY: clean
clean:
	-rm -f zebra
	-rm -f zebra-server
	-rm -f coverage.txt
	-rm -rf ./bin
	-rm -f ./simulator/*.crt
	-rm -f ./simulator/*.key
	-rm -rf ./simulator/simulator-store
	-rm -rf ./simulator/*.crt
	-rm -rf ./simulator/*.csr
	-rm -rf ./simulator/*.key
	-rm -rf ./simulator/*.json
	-rm -rf ./simulator/*.yaml
