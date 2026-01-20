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
	sol-fmt \
	sol-lint \
	sol-sec \
	sol-test-unit \
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

################################################################################
# Solidity Targets
################################################################################
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
	@slither ./contracts/src

.PHONY: sol-test-unit
sol-test-unit: sol-build
	@forge test -vvv

################################################################################
# Shared Targets
################################################################################
process_compose_port=8079
process_compose_config=.process-compose.yaml
integration_dir=./test/integration

.PHONY: test-integration
test-integration: \
	test-integration-hello-world

.PHONY: test-integration-hello-world
test-integration-hello-world: sol-build tidy
	@process-compose up \
		--tui=false \
		--port=${process_compose_port} \
		-f ${integration_dir}/hello-world/${process_compose_config} \
		2> /dev/null

.PHONY: clean
clean:
	@forge clean
