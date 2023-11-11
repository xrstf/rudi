.PHONY: generate
generate:
	mkdir -p pkg/lang/parser
	pigeon corel.peg > pkg/lang/parser/generated.go

.PHONY: clean
clean:
	rm -rf _build

.PHONY: build
build:
	mkdir -p _build
	go build -v -o _build/ ./cmd/tester

.PHONY: test
test:
	_build/tester test.corel
