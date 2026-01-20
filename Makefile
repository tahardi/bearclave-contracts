# https://clarkgrubb.com/makefile-style-guide
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := pre-pr
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: pre-pr
pre-pr: tidy sol-fmt sol-lint sol-sec sol-test

# https://golangci-lint.run/welcome/install/#install-from-sources
# They do not recommend using golangci-lint via go tool directive
# as there are still bugs, but I want to try out go tool and work
# uses an old version of golangci-lint. So, I don't mind guinea
# pigging go tool and using a new version of golangci-lint in here
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

.PHONY: sol-test
sol-test: sol-build
	@forge test -vvv

.PHONY: clean
clean:
	@forge clean
