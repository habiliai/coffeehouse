SHELL := $(shell which sh)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
PROTOC_DIR := bin/protoc-$(GOOS)-$(GOARCH)
PROTOC := bin/protoc

.PHONY: all
all: bin/alice

bin/protoc-linux-amd64.zip:
	wget -O $@ "https://github.com/protocolbuffers/protobuf/releases/download/v27.1/protoc-27.1-linux-x86_64.zip"

bin/protoc-linux-arm64.zip:
	wget -O $@ "https://github.com/protocolbuffers/protobuf/releases/download/v27.1/protoc-27.1-linux-aarch_64.zip"

bin/protoc-darwin-amd64.zip:
	wget -O $@ "https://github.com/protocolbuffers/protobuf/releases/download/v27.1/protoc-27.1-osx-x86_64.zip"

bin/protoc-darwin-arm64.zip:
	wget -O $@ "https://github.com/protocolbuffers/protobuf/releases/download/v27.1/protoc-27.1-osx-aarch_64.zip"

$(PROTOC_DIR): bin/protoc-$(GOOS)-$(GOARCH).zip
	@unzip -o $< -d $(PROTOC_DIR)
	@echo "done $@"

$(PROTOC): $(PROTOC_DIR)
	chmod 755 $(PROTOC_DIR)/bin/protoc
	ln -sf protoc-$(GOOS)-$(GOARCH)/bin/protoc $(PROTOC)
	touch $(PROTOC)

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
$(PROTOC_GEN_GO):
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34

PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
$(PROTOC_GEN_GO_GRPC):
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4

%.pb.go: %.proto $(PROTOC) $(PROTOC_GEN_GO)
	@export PATH="$(shell go env GOPATH)/bin:$(PATH)"
	$(PROTOC) --go_out=. --go_opt=paths=source_relative -I. $<

%_grpc.pb.go: %.proto $(PROTOC) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)
	@export PATH="$(shell go env GOPATH)/bin:$(PATH)"
	$(PROTOC) --go-grpc_out=. --go-grpc_opt=paths=source_relative -I. $<

pb: aliceapi/aliceapi_grpc.pb.go aliceapi/aliceapi.pb.go

./bin/golangci-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.62.2
	@chmod +x ./bin/golangci-lint
	@echo "golangci-lint installed"

.PHONY: lint
lint: ./bin/golangci-lint $(PB_FILES)
	./bin/golangci-lint run

.PHONY: test
test: pb
	PROJECT_DIR=$(shell pwd) CI=true go test -timeout 15m -p 1 ./...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf $(PB_FILES)
	rm -rf $(PROTOC_DIR)
	rm -f $(PROTOC)
	@echo "cleared"

.PHONY: bin/alice
bin/alice: pb
	CGO_ENABLED=0 go build -o $@ ./cmd/alice

.PHONY: seed
seed: bin/alice
	bin/alice seed --reset
