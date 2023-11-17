# SPDX-FileCopyrightText: 2023 Christoph Mewes
# SPDX-License-Identifier: MIT

GIT_VERSION = $(shell git describe --tags --always)
GIT_HEAD ?= $(shell git log -1 --format=%H)
NOW_GO_RFC339 = $(shell date --utc +'%Y-%m-%dT%H:%M:%SZ')

export CGO_ENABLED ?= 0
export GOFLAGS ?= -mod=readonly -trimpath
OUTPUT_DIR ?= _build
GO_DEFINES ?= -X main.BuildTag=$(GIT_VERSION) -X main.BuildCommit=$(GIT_HEAD) -X main.BuildDate=$(NOW_GO_RFC339)
GO_LDFLAGS += -w -extldflags '-static' $(GO_DEFINES)
GO_BUILD_FLAGS ?= -v -ldflags '$(GO_LDFLAGS)'
GO_TEST_FLAGS ?=

.PHONY: generate
generate:
	pigeon pkg/lang/grammar/otto.peg > pkg/lang/parser/generated.go

.PHONY: clean
clean:
	rm -rf _build

.PHONY: build
build:
	mkdir -p _build
	go build $(GO_BUILD_FLAGS) -o _build/ ./cmd/tester
	cd cmd/otti && go build $(GO_BUILD_FLAGS) -o ../../_build .

.PHONY: test
test:
	CGO_ENABLED=1 go test $(GO_TEST_FLAGS) ./...

.PHONY: lint
lint:
	golangci-lint run ./...
