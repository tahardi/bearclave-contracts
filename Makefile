# https://clarkgrubb.com/makefile-style-guide
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := pre-pr
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: pre-pr
pre-pr: \
	tidy \
	mock \
	go-lint \
	go-test \
	sol-fmt \
	sol-lint \
	sol-sec \
	sol-test \
	bindings \
	test-integration

################################################################################
# Go Targets
################################################################################
lint_modfile=modfiles/golangci-lint/go.mod
.PHONY: go-lint
go-lint:
	@go tool -modfile=$(lint_modfile) golangci-lint run --config .golangci.yaml

.PHONY: go-lint-version
go-lint-version:
	@go tool -modfile=$(lint_modfile) golangci-lint --version

mockery_modfile=modfiles/mockery/go.mod
.PHONY: mock
mock: tidy
	@go tool -modfile=$(mockery_modfile) mockery --config=.mockery.yaml

.PHONY: mock-version
mock-version:
	@go tool -modfile=$(mockery_modfile) mockery version

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: go-test
go-test: go-test-foundry

.PHONY: go-test-foundry
go-test-foundry:
	@go test -v -count=1 -race ./test/foundry/...

################################################################################
# Solidity Targets
################################################################################
contracts_dir=./contracts
broadcast_dir=$(contracts_dir)/broadcast
cache_dir=$(contracts_dir)/cache
out_dir=$(contracts_dir)/out
src_dir=$(contracts_dir)/src

.PHONY: sol-build
sol-build:
	@forge build

.PHONY: sol-fmt
sol-fmt:
	@forge fmt

.PHONY: sol-lint
sol-lint:
	@forge lint

.PHONY: sol-sec
sol-sec:
	@slither $(src_dir) --config-file .slither-config.json

.PHONY: sol-test
sol-test: sol-build
	@forge test -vvv

################################################################################
# Shared Targets
################################################################################
bindings_dir=$(contracts_dir)/bindings
bindings_pkg=bindings
integration_dir=./test/integration

.PHONY: bindings
bindings: \
	bindings-bear-coin \
	bindings-hello-world

.PHONY: bindings-bear-coin
bindings-bear-coin: sol-build
	@jq '.abi' $(out_dir)/BearCoin.sol/BearCoin.json | \
	abigen \
		--abi /dev/stdin \
		--pkg $(bindings_pkg) \
		--type BearCoin \
		--out $(bindings_dir)/bearcoin.go

.PHONY: bindings-hello-world
bindings-hello-world: sol-build
	@jq '.abi' $(out_dir)/HelloWorld.sol/HelloWorld.json | \
	abigen \
		--abi /dev/stdin \
		--pkg $(bindings_pkg) \
		--type HelloWorld \
		--out $(bindings_dir)/helloworld.go

.PHONY: test-integration
test-integration: \
	test-integration-hello-world

.PHONY: test-integration-hello-world
test-integration-hello-world: sol-build tidy
	@go test -v -count=1 $(integration_dir)/hello-world/helloworld_test.go

.PHONY: clean
clean:
	forge clean
	rm -rf $(broadcast_dir)
	rm -rf $(cache_dir)
