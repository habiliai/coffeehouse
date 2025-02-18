SHELL := $(shell which sh)
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
AF_PB_DIR := ./pkg/habapi
AF_PB_FILES := $(AF_PB_DIR)/habapi.pb.go
PROTO_DIR := ../proto
AF_PROTO_FILES := $(PROTO_DIR)/habapi.proto
PROTOC_DIR := bin/protoc-$(GOOS)-$(GOARCH)
PROTOC := bin/protoc

.PHONY: all
all: bin/habapi

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

$(AF_PB_FILES): $(AF_PROTO_FILES) $(PROTOC)
	$(eval PROTOC_ABS := $(abspath $(PROTOC)))
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4
	@export PATH="$(shell go env GOPATH)/bin:$(PATH)"
	$(PROTOC_ABS) --go_out=$(AF_PB_DIR) --go_opt=paths=source_relative --go-grpc_out=$(AF_PB_DIR) --go-grpc_opt=paths=source_relative -I$(PROTO_DIR) $(AF_PROTO_FILES)

pb: $(AF_PB_FILES)

./bin/golangci-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.62.2
	@chmod +x ./bin/golangci-lint
	@echo "golangci-lint installed"

.PHONY: lint
lint: ./bin/golangci-lint $(PB_FILES)
	./bin/golangci-lint run

.PHONY: test
test: pb
	CI=true go test -timeout 15m -p 1 ./...

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf $(PB_FILES)
	rm -rf $(PROTOC_DIR)
	rm -f $(PROTOC)
	@echo "cleared"

bin/habapi: pb
	CGO_ENABLED=0 go build -o bin/habapi ./cmd/habapi/main.go

.PHONY: seed
seed: bin/habapi
	bin/habapi seed --reset
