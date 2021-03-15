SHELL := /bin/bash

all: \
	buf-lint \
	buf-generate \
	go-lint \
	go-review \
	buf-generate \
	eslint \
	go-test \
	go-mod-tidy \
	git-verify-nodiff

include tools/buf/rules.mk
include tools/commitlint/rules.mk
include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/goreview/rules.mk
include tools/semantic-release/rules.mk
include tools/eslint/rules.mk

.PHONY: examples/proto/api-common-protos
examples/proto/api-common-protos:
	@git submodule update --init --recursive $@

.PHONY: go-test
go-test:
	$(info [$@] running Go tests...)
	@go test -count 1 -cover -race ./...

.PHONY: go-mod-tidy
go-mod-tidy:
	$(info [$@] tidying Go module files...)
	@go mod tidy -v

.PHONY: buf-lint
buf-lint: $(buf) examples/proto/api-common-protos
	$(info [$@] linting protobuf schemas...)
	@$(buf) lint

protoc_gen_typescript_http := ./bin/protoc-gen-typescript-http
export PATH := $(dir $(abspath $(protoc_gen_typescript_http))):$(PATH)

.PHONY: $(protoc_gen_typescript_http)
$(protoc_gen_typescript_http):
	$(info [$@] building protoc-gen-typescript-http...)
	@go build -o $@ .

.PHONY: $(eslint)
eslint: $(eslint)
	$(info [$@] linting typescript files...)
	$(eslint) --config $(eslint_cwd)/.eslintrc.js --quiet "examples/proto/gen/typescript/**/*.ts"

.PHONY: buf-generate
buf-generate: $(buf) $(protoc_gen_typescript_http) examples/proto/api-common-protos
	$(info [$@] generating protobuf stubs...)
	@rm -rf examples/proto/gen
	@$(buf) generate --path examples/proto/src/einride
