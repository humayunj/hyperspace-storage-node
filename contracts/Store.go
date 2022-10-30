// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package store

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
)

// StoreMetaData contains all meta data concerning the Store contract.
var StoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_TLSCert\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"_HOST\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"HOST\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TLSCert\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumStorageNode.CallerType\",\"name\":\"callerType\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"},{\"internalType\":\"bytes20\",\"name\":\"merkleRootHash\",\"type\":\"bytes20\"},{\"internalType\":\"uint32\",\"name\":\"fileSize\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"timerStart\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timerEnd\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"proveTimeoutLength\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"concludeTimeoutLength\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"segmentsCount\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"bidAmount\",\"type\":\"uint256\"}],\"name\":\"concludeTransaction\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"finishTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockedCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mappingLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"mappingsList\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"processProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validateStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// StoreABI is the input ABI used to generate the binding from.
// Deprecated: Use StoreMetaData.ABI instead.
var StoreABI = StoreMetaData.ABI

// Store is an auto generated Go binding around an Ethereum contract.
type Store struct {
	StoreCaller     // Read-only binding to the contract
	StoreTransactor // Write-only binding to the contract
	StoreFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type StoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StoreSession struct {
	Contract     *Store            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StoreCallerSession struct {
	Contract *StoreCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StoreTransactorSession struct {
	Contract     *StoreTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type StoreRaw struct {
	Contract *Store // Generic contract binding to access the raw methods on
}

// StoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StoreCallerRaw struct {
	Contract *StoreCaller // Generic read-only contract binding to access the raw methods on
}

// StoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StoreTransactorRaw struct {
	Contract *StoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStore creates a new instance of Store, bound to a specific deployed contract.
func NewStore(address common.Address, backend bind.ContractBackend) (*Store, error) {
	contract, err := bindStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Store{StoreCaller: StoreCaller{contract: contract}, StoreTransactor: StoreTransactor{contract: contract}, StoreFilterer: StoreFilterer{contract: contract}}, nil
}

// NewStoreCaller creates a new read-only instance of Store, bound to a specific deployed contract.
func NewStoreCaller(address common.Address, caller bind.ContractCaller) (*StoreCaller, error) {
	contract, err := bindStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StoreCaller{contract: contract}, nil
}

// NewStoreTransactor creates a new write-only instance of Store, bound to a specific deployed contract.
func NewStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*StoreTransactor, error) {
	contract, err := bindStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StoreTransactor{contract: contract}, nil
}

// NewStoreFilterer creates a new log filterer instance of Store, bound to a specific deployed contract.
func NewStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*StoreFilterer, error) {
	contract, err := bindStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StoreFilterer{contract: contract}, nil
}

// bindStore binds a generic wrapper to an already deployed contract.
func bindStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.StoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Store.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.contract.Transact(opts, method, params...)
}

// HOST is a free data retrieval call binding the contract method 0x49f289dc.
//
// Solidity: function HOST() view returns(string)
func (_Store *StoreCaller) HOST(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "HOST")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HOST is a free data retrieval call binding the contract method 0x49f289dc.
//
// Solidity: function HOST() view returns(string)
func (_Store *StoreSession) HOST() (string, error) {
	return _Store.Contract.HOST(&_Store.CallOpts)
}

// HOST is a free data retrieval call binding the contract method 0x49f289dc.
//
// Solidity: function HOST() view returns(string)
func (_Store *StoreCallerSession) HOST() (string, error) {
	return _Store.Contract.HOST(&_Store.CallOpts)
}

// TLSCert is a free data retrieval call binding the contract method 0xb3752fa3.
//
// Solidity: function TLSCert() view returns(bytes)
func (_Store *StoreCaller) TLSCert(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "TLSCert")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TLSCert is a free data retrieval call binding the contract method 0xb3752fa3.
//
// Solidity: function TLSCert() view returns(bytes)
func (_Store *StoreSession) TLSCert() ([]byte, error) {
	return _Store.Contract.TLSCert(&_Store.CallOpts)
}

// TLSCert is a free data retrieval call binding the contract method 0xb3752fa3.
//
// Solidity: function TLSCert() view returns(bytes)
func (_Store *StoreCallerSession) TLSCert() ([]byte, error) {
	return _Store.Contract.TLSCert(&_Store.CallOpts)
}

// LockedCollateral is a free data retrieval call binding the contract method 0xb952cc4a.
//
// Solidity: function lockedCollateral() view returns(uint256)
func (_Store *StoreCaller) LockedCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "lockedCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedCollateral is a free data retrieval call binding the contract method 0xb952cc4a.
//
// Solidity: function lockedCollateral() view returns(uint256)
func (_Store *StoreSession) LockedCollateral() (*big.Int, error) {
	return _Store.Contract.LockedCollateral(&_Store.CallOpts)
}

// LockedCollateral is a free data retrieval call binding the contract method 0xb952cc4a.
//
// Solidity: function lockedCollateral() view returns(uint256)
func (_Store *StoreCallerSession) LockedCollateral() (*big.Int, error) {
	return _Store.Contract.LockedCollateral(&_Store.CallOpts)
}

// MappingLength is a free data retrieval call binding the contract method 0x116766a6.
//
// Solidity: function mappingLength() view returns(uint256)
func (_Store *StoreCaller) MappingLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "mappingLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MappingLength is a free data retrieval call binding the contract method 0x116766a6.
//
// Solidity: function mappingLength() view returns(uint256)
func (_Store *StoreSession) MappingLength() (*big.Int, error) {
	return _Store.Contract.MappingLength(&_Store.CallOpts)
}

// MappingLength is a free data retrieval call binding the contract method 0x116766a6.
//
// Solidity: function mappingLength() view returns(uint256)
func (_Store *StoreCallerSession) MappingLength() (*big.Int, error) {
	return _Store.Contract.MappingLength(&_Store.CallOpts)
}

// MappingsList is a free data retrieval call binding the contract method 0xca88afbc.
//
// Solidity: function mappingsList(uint256 ) view returns(address)
func (_Store *StoreCaller) MappingsList(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Store.contract.Call(opts, &out, "mappingsList", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MappingsList is a free data retrieval call binding the contract method 0xca88afbc.
//
// Solidity: function mappingsList(uint256 ) view returns(address)
func (_Store *StoreSession) MappingsList(arg0 *big.Int) (common.Address, error) {
	return _Store.Contract.MappingsList(&_Store.CallOpts, arg0)
}

// MappingsList is a free data retrieval call binding the contract method 0xca88afbc.
//
// Solidity: function mappingsList(uint256 ) view returns(address)
func (_Store *StoreCallerSession) MappingsList(arg0 *big.Int) (common.Address, error) {
	return _Store.Contract.MappingsList(&_Store.CallOpts, arg0)
}

// ConcludeTransaction is a paid mutator transaction binding the contract method 0x15e23f6b.
//
// Solidity: function concludeTransaction(uint8 callerType, address userAddress, bytes20 merkleRootHash, uint32 fileSize, uint256 timerStart, uint256 timerEnd, uint64 proveTimeoutLength, uint64 concludeTimeoutLength, uint32 segmentsCount, uint256 bidAmount) payable returns()
func (_Store *StoreTransactor) ConcludeTransaction(opts *bind.TransactOpts, callerType uint8, userAddress common.Address, merkleRootHash [20]byte, fileSize uint32, timerStart *big.Int, timerEnd *big.Int, proveTimeoutLength uint64, concludeTimeoutLength uint64, segmentsCount uint32, bidAmount *big.Int) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "concludeTransaction", callerType, userAddress, merkleRootHash, fileSize, timerStart, timerEnd, proveTimeoutLength, concludeTimeoutLength, segmentsCount, bidAmount)
}

// ConcludeTransaction is a paid mutator transaction binding the contract method 0x15e23f6b.
//
// Solidity: function concludeTransaction(uint8 callerType, address userAddress, bytes20 merkleRootHash, uint32 fileSize, uint256 timerStart, uint256 timerEnd, uint64 proveTimeoutLength, uint64 concludeTimeoutLength, uint32 segmentsCount, uint256 bidAmount) payable returns()
func (_Store *StoreSession) ConcludeTransaction(callerType uint8, userAddress common.Address, merkleRootHash [20]byte, fileSize uint32, timerStart *big.Int, timerEnd *big.Int, proveTimeoutLength uint64, concludeTimeoutLength uint64, segmentsCount uint32, bidAmount *big.Int) (*types.Transaction, error) {
	return _Store.Contract.ConcludeTransaction(&_Store.TransactOpts, callerType, userAddress, merkleRootHash, fileSize, timerStart, timerEnd, proveTimeoutLength, concludeTimeoutLength, segmentsCount, bidAmount)
}

// ConcludeTransaction is a paid mutator transaction binding the contract method 0x15e23f6b.
//
// Solidity: function concludeTransaction(uint8 callerType, address userAddress, bytes20 merkleRootHash, uint32 fileSize, uint256 timerStart, uint256 timerEnd, uint64 proveTimeoutLength, uint64 concludeTimeoutLength, uint32 segmentsCount, uint256 bidAmount) payable returns()
func (_Store *StoreTransactorSession) ConcludeTransaction(callerType uint8, userAddress common.Address, merkleRootHash [20]byte, fileSize uint32, timerStart *big.Int, timerEnd *big.Int, proveTimeoutLength uint64, concludeTimeoutLength uint64, segmentsCount uint32, bidAmount *big.Int) (*types.Transaction, error) {
	return _Store.Contract.ConcludeTransaction(&_Store.TransactOpts, callerType, userAddress, merkleRootHash, fileSize, timerStart, timerEnd, proveTimeoutLength, concludeTimeoutLength, segmentsCount, bidAmount)
}

// FinishTransaction is a paid mutator transaction binding the contract method 0x34cc9e41.
//
// Solidity: function finishTransaction(address userAddress) returns()
func (_Store *StoreTransactor) FinishTransaction(opts *bind.TransactOpts, userAddress common.Address) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "finishTransaction", userAddress)
}

// FinishTransaction is a paid mutator transaction binding the contract method 0x34cc9e41.
//
// Solidity: function finishTransaction(address userAddress) returns()
func (_Store *StoreSession) FinishTransaction(userAddress common.Address) (*types.Transaction, error) {
	return _Store.Contract.FinishTransaction(&_Store.TransactOpts, userAddress)
}

// FinishTransaction is a paid mutator transaction binding the contract method 0x34cc9e41.
//
// Solidity: function finishTransaction(address userAddress) returns()
func (_Store *StoreTransactorSession) FinishTransaction(userAddress common.Address) (*types.Transaction, error) {
	return _Store.Contract.FinishTransaction(&_Store.TransactOpts, userAddress)
}

// ProcessProof is a paid mutator transaction binding the contract method 0xcb89c46b.
//
// Solidity: function processProof() returns()
func (_Store *StoreTransactor) ProcessProof(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "processProof")
}

// ProcessProof is a paid mutator transaction binding the contract method 0xcb89c46b.
//
// Solidity: function processProof() returns()
func (_Store *StoreSession) ProcessProof() (*types.Transaction, error) {
	return _Store.Contract.ProcessProof(&_Store.TransactOpts)
}

// ProcessProof is a paid mutator transaction binding the contract method 0xcb89c46b.
//
// Solidity: function processProof() returns()
func (_Store *StoreTransactorSession) ProcessProof() (*types.Transaction, error) {
	return _Store.Contract.ProcessProof(&_Store.TransactOpts)
}

// ValidateStorage is a paid mutator transaction binding the contract method 0x599b2e31.
//
// Solidity: function validateStorage() returns()
func (_Store *StoreTransactor) ValidateStorage(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "validateStorage")
}

// ValidateStorage is a paid mutator transaction binding the contract method 0x599b2e31.
//
// Solidity: function validateStorage() returns()
func (_Store *StoreSession) ValidateStorage() (*types.Transaction, error) {
	return _Store.Contract.ValidateStorage(&_Store.TransactOpts)
}

// ValidateStorage is a paid mutator transaction binding the contract method 0x599b2e31.
//
// Solidity: function validateStorage() returns()
func (_Store *StoreTransactorSession) ValidateStorage() (*types.Transaction, error) {
	return _Store.Contract.ValidateStorage(&_Store.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x00f714ce.
//
// Solidity: function withdraw(uint256 amount, address target) returns()
func (_Store *StoreTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int, target common.Address) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "withdraw", amount, target)
}

// Withdraw is a paid mutator transaction binding the contract method 0x00f714ce.
//
// Solidity: function withdraw(uint256 amount, address target) returns()
func (_Store *StoreSession) Withdraw(amount *big.Int, target common.Address) (*types.Transaction, error) {
	return _Store.Contract.Withdraw(&_Store.TransactOpts, amount, target)
}

// Withdraw is a paid mutator transaction binding the contract method 0x00f714ce.
//
// Solidity: function withdraw(uint256 amount, address target) returns()
func (_Store *StoreTransactorSession) Withdraw(amount *big.Int, target common.Address) (*types.Transaction, error) {
	return _Store.Contract.Withdraw(&_Store.TransactOpts, amount, target)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Store *StoreTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Store *StoreSession) Receive() (*types.Transaction, error) {
	return _Store.Contract.Receive(&_Store.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Store *StoreTransactorSession) Receive() (*types.Transaction, error) {
	return _Store.Contract.Receive(&_Store.TransactOpts)
}
