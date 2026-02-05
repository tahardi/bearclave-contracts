# Bearclave: Smart Contracts

This repository was created to explore Ethereum-based blockchain development
tools and applications. Eventually, the work here will tie into the
[Bearclave](https://github.com/tahardi/bearclave)
project through end-to-end Data Oracle and Off-Chain Compute examples.

## Getting Started

Follow the steps below to install and setup the tools necessary for compiling,
deploying, and testing Bearclave smart contracts and go bindings.

1. Install [Golang](https://golang.org/doc/install) (v1.25.5 or higher) to build
   and run the integration tests.
2. Install the `abigen` tool to generate Go bindings for smart contracts.
```bash
 go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```
3. Install the [Foundry](https://github.com/foundry-rs/foundry) toolset for
   smart contract development.
4. Initialize the Forge and OpenZeppelin submodules.
```bash
git submodule update --init --recursive 
```
5. Install [Slither](https://github.com/crytic/slither) static analysis tool for
   auditing smart contracts.
