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

.PHONY: all
all: clean generate docs build test spellcheck

.PHONY: generate
generate:
	pigeon pkg/lang/grammar/rudi.peg > pkg/lang/parser/generated.go

.PHONY: docs
docs:
	cd hack/docs-toc && go mod tidy && go run .
	cd hack/docs-prerender && go mod tidy && go run .

.PHONY: clean
clean:
	rm -rf _build

.PHONY: build
build:
	mkdir -p _build
	cd cmd/rudi && go build $(GO_BUILD_FLAGS) -o ../../_build .
	cd cmd/example && go build $(GO_BUILD_FLAGS) -o ../../_build .

.PHONY: install
install:
	cd cmd/rudi && go install $(GO_BUILD_FLAGS) .

.PHONY: test
test:
	CGO_ENABLED=1 go test $(GO_TEST_FLAGS) ./...

.PHONY: lint
lint:
	golangci-lint run ./...
	cd cmd/rudi && golangci-lint run ./... --config ../../.golangci.yml
	cd hack/docs-toc && golangci-lint run ./... --config ../../.golangci.yml
	cd hack/docs-prerender && golangci-lint run ./... --config ../../.golangci.yml

.PHONY: spellcheck
spellcheck:
	typos
