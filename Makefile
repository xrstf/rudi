.PHONY: generate
generate:
	pigeon grammar/otto.peg > pkg/lang/parser/generated.go

.PHONY: clean
clean:
	rm -rf _build

.PHONY: build
build:
	mkdir -p _build
	go build -v -o _build/ ./cmd/tester

.PHONY: run-tests
run-tests:
	go test -v ./...

.PHONY: test
test:
	_build/tester test.otto
