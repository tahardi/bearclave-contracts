// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// HelloWorldMetaData contains all meta data concerning the HelloWorld contract.
var HelloWorldMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"greet\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"}]",
}

// HelloWorldABI is the input ABI used to generate the binding from.
// Deprecated: Use HelloWorldMetaData.ABI instead.
var HelloWorldABI = HelloWorldMetaData.ABI

// HelloWorld is an auto generated Go binding around an Ethereum contract.
type HelloWorld struct {
	HelloWorldCaller     // Read-only binding to the contract
	HelloWorldTransactor // Write-only binding to the contract
	HelloWorldFilterer   // Log filterer for contract events
}

// HelloWorldCaller is an auto generated read-only Go binding around an Ethereum contract.
type HelloWorldCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelloWorldTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HelloWorldTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelloWorldFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HelloWorldFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HelloWorldSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HelloWorldSession struct {
	Contract     *HelloWorld       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HelloWorldCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HelloWorldCallerSession struct {
	Contract *HelloWorldCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// HelloWorldTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HelloWorldTransactorSession struct {
	Contract     *HelloWorldTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// HelloWorldRaw is an auto generated low-level Go binding around an Ethereum contract.
type HelloWorldRaw struct {
	Contract *HelloWorld // Generic contract binding to access the raw methods on
}

// HelloWorldCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HelloWorldCallerRaw struct {
	Contract *HelloWorldCaller // Generic read-only contract binding to access the raw methods on
}

// HelloWorldTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HelloWorldTransactorRaw struct {
	Contract *HelloWorldTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHelloWorld creates a new instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorld(address common.Address, backend bind.ContractBackend) (*HelloWorld, error) {
	contract, err := bindHelloWorld(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HelloWorld{HelloWorldCaller: HelloWorldCaller{contract: contract}, HelloWorldTransactor: HelloWorldTransactor{contract: contract}, HelloWorldFilterer: HelloWorldFilterer{contract: contract}}, nil
}

// NewHelloWorldCaller creates a new read-only instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorldCaller(address common.Address, caller bind.ContractCaller) (*HelloWorldCaller, error) {
	contract, err := bindHelloWorld(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HelloWorldCaller{contract: contract}, nil
}

// NewHelloWorldTransactor creates a new write-only instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorldTransactor(address common.Address, transactor bind.ContractTransactor) (*HelloWorldTransactor, error) {
	contract, err := bindHelloWorld(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HelloWorldTransactor{contract: contract}, nil
}

// NewHelloWorldFilterer creates a new log filterer instance of HelloWorld, bound to a specific deployed contract.
func NewHelloWorldFilterer(address common.Address, filterer bind.ContractFilterer) (*HelloWorldFilterer, error) {
	contract, err := bindHelloWorld(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HelloWorldFilterer{contract: contract}, nil
}

// bindHelloWorld binds a generic wrapper to an already deployed contract.
func bindHelloWorld(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HelloWorldMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HelloWorld *HelloWorldRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HelloWorld.Contract.HelloWorldCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HelloWorld *HelloWorldRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HelloWorld.Contract.HelloWorldTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HelloWorld *HelloWorldRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HelloWorld.Contract.HelloWorldTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HelloWorld *HelloWorldCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HelloWorld.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HelloWorld *HelloWorldTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HelloWorld.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HelloWorld *HelloWorldTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HelloWorld.Contract.contract.Transact(opts, method, params...)
}

// Greet is a free data retrieval call binding the contract method 0xcfae3217.
//
// Solidity: function greet() view returns(string)
func (_HelloWorld *HelloWorldCaller) Greet(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HelloWorld.contract.Call(opts, &out, "greet")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Greet is a free data retrieval call binding the contract method 0xcfae3217.
//
// Solidity: function greet() view returns(string)
func (_HelloWorld *HelloWorldSession) Greet() (string, error) {
	return _HelloWorld.Contract.Greet(&_HelloWorld.CallOpts)
}

// Greet is a free data retrieval call binding the contract method 0xcfae3217.
//
// Solidity: function greet() view returns(string)
func (_HelloWorld *HelloWorldCallerSession) Greet() (string, error) {
	return _HelloWorld.Contract.Greet(&_HelloWorld.CallOpts)
}
