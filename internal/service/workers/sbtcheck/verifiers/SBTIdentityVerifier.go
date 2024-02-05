// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package verifiers

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

// IBaseVerifierProveIdentityParams is an auto generated low-level Go binding around an user-defined struct.
type IBaseVerifierProveIdentityParams struct {
	StatesMerkleData ILightweightStateStatesMerkleData
	Inputs           []*big.Int
	A                [2]*big.Int
	B                [2][2]*big.Int
	C                [2]*big.Int
}

// IBaseVerifierTransitStateParams is an auto generated low-level Go binding around an user-defined struct.
type IBaseVerifierTransitStateParams struct {
	NewIdentitiesStatesRoot [32]byte
	GistData                ILightweightStateGistRootData
	Proof                   []byte
}

// ILightweightStateGistRootData is an auto generated low-level Go binding around an user-defined struct.
type ILightweightStateGistRootData struct {
	Root               *big.Int
	CreatedAtTimestamp *big.Int
}

// ILightweightStateStatesMerkleData is an auto generated low-level Go binding around an user-defined struct.
type ILightweightStateStatesMerkleData struct {
	IssuerId           *big.Int
	IssuerState        *big.Int
	CreatedAtTimestamp *big.Int
	MerkleProof        [][32]byte
}

// ISBTIdentityVerifierSBTIdentityProofInfo is an auto generated low-level Go binding around an user-defined struct.
type ISBTIdentityVerifierSBTIdentityProofInfo struct {
	SenderAddr common.Address
	SbtTokenId *big.Int
	IsProved   bool
}

// SBTIdentityVerifierMetaData contains all meta data concerning the SBTIdentityVerifier contract.
var SBTIdentityVerifierMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"identityId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"SBTIdentityProved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"SBT_IDENTITY_PROOF_QUERY_ID\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIZKPQueriesStorage\",\"name\":\"zkpQueriesStorage_\",\"type\":\"address\"},{\"internalType\":\"contractIVerifiedSBT\",\"name\":\"sbtToken_\",\"type\":\"address\"}],\"name\":\"__SBTIdentityVerifier_init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addressToIdentityId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"schema_\",\"type\":\"uint256\"}],\"name\":\"getAllowedIssuers\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"identityId_\",\"type\":\"uint256\"}],\"name\":\"getIdentityProofInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"senderAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sbtTokenId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isProved\",\"type\":\"bool\"}],\"internalType\":\"structISBTIdentityVerifier.SBTIdentityProofInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"schema_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issuerId_\",\"type\":\"uint256\"}],\"name\":\"isAllowedIssuer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddr_\",\"type\":\"address\"}],\"name\":\"isIdentityProved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"identityId_\",\"type\":\"uint256\"}],\"name\":\"isIdentityProved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"issuerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issuerState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"internalType\":\"structILightweightState.StatesMerkleData\",\"name\":\"statesMerkleData\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"inputs\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"internalType\":\"structIBaseVerifier.ProveIdentityParams\",\"name\":\"proveIdentityParams_\",\"type\":\"tuple\"}],\"name\":\"proveIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sbtToken\",\"outputs\":[{\"internalType\":\"contractIVerifiedSBT\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIZKPQueriesStorage\",\"name\":\"newZKPQueriesStorage_\",\"type\":\"address\"}],\"name\":\"setZKPQueriesStorage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"issuerId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"issuerState\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"merkleProof\",\"type\":\"bytes32[]\"}],\"internalType\":\"structILightweightState.StatesMerkleData\",\"name\":\"statesMerkleData\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"inputs\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"internalType\":\"structIBaseVerifier.ProveIdentityParams\",\"name\":\"proveIdentityParams_\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"newIdentitiesStatesRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"root\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAtTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structILightweightState.GistRootData\",\"name\":\"gistData\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"}],\"internalType\":\"structIBaseVerifier.TransitStateParams\",\"name\":\"transitStateParams_\",\"type\":\"tuple\"}],\"name\":\"transitStateAndProveIdentity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"schema_\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"issuerIds_\",\"type\":\"uint256[]\"},{\"internalType\":\"bool\",\"name\":\"isAdding_\",\"type\":\"bool\"}],\"name\":\"updateAllowedIssuers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zkpQueriesStorage\",\"outputs\":[{\"internalType\":\"contractIZKPQueriesStorage\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SBTIdentityVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use SBTIdentityVerifierMetaData.ABI instead.
var SBTIdentityVerifierABI = SBTIdentityVerifierMetaData.ABI

// SBTIdentityVerifier is an auto generated Go binding around an Ethereum contract.
type SBTIdentityVerifier struct {
	SBTIdentityVerifierCaller     // Read-only binding to the contract
	SBTIdentityVerifierTransactor // Write-only binding to the contract
	SBTIdentityVerifierFilterer   // Log filterer for contract events
}

// SBTIdentityVerifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type SBTIdentityVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SBTIdentityVerifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SBTIdentityVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SBTIdentityVerifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SBTIdentityVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SBTIdentityVerifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SBTIdentityVerifierSession struct {
	Contract     *SBTIdentityVerifier // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// SBTIdentityVerifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SBTIdentityVerifierCallerSession struct {
	Contract *SBTIdentityVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// SBTIdentityVerifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SBTIdentityVerifierTransactorSession struct {
	Contract     *SBTIdentityVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// SBTIdentityVerifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type SBTIdentityVerifierRaw struct {
	Contract *SBTIdentityVerifier // Generic contract binding to access the raw methods on
}

// SBTIdentityVerifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SBTIdentityVerifierCallerRaw struct {
	Contract *SBTIdentityVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// SBTIdentityVerifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SBTIdentityVerifierTransactorRaw struct {
	Contract *SBTIdentityVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSBTIdentityVerifier creates a new instance of SBTIdentityVerifier, bound to a specific deployed contract.
func NewSBTIdentityVerifier(address common.Address, backend bind.ContractBackend) (*SBTIdentityVerifier, error) {
	contract, err := bindSBTIdentityVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifier{SBTIdentityVerifierCaller: SBTIdentityVerifierCaller{contract: contract}, SBTIdentityVerifierTransactor: SBTIdentityVerifierTransactor{contract: contract}, SBTIdentityVerifierFilterer: SBTIdentityVerifierFilterer{contract: contract}}, nil
}

// NewSBTIdentityVerifierCaller creates a new read-only instance of SBTIdentityVerifier, bound to a specific deployed contract.
func NewSBTIdentityVerifierCaller(address common.Address, caller bind.ContractCaller) (*SBTIdentityVerifierCaller, error) {
	contract, err := bindSBTIdentityVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierCaller{contract: contract}, nil
}

// NewSBTIdentityVerifierTransactor creates a new write-only instance of SBTIdentityVerifier, bound to a specific deployed contract.
func NewSBTIdentityVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*SBTIdentityVerifierTransactor, error) {
	contract, err := bindSBTIdentityVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierTransactor{contract: contract}, nil
}

// NewSBTIdentityVerifierFilterer creates a new log filterer instance of SBTIdentityVerifier, bound to a specific deployed contract.
func NewSBTIdentityVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*SBTIdentityVerifierFilterer, error) {
	contract, err := bindSBTIdentityVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierFilterer{contract: contract}, nil
}

// bindSBTIdentityVerifier binds a generic wrapper to an already deployed contract.
func bindSBTIdentityVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SBTIdentityVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SBTIdentityVerifier *SBTIdentityVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SBTIdentityVerifier.Contract.SBTIdentityVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SBTIdentityVerifier *SBTIdentityVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.SBTIdentityVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SBTIdentityVerifier *SBTIdentityVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.SBTIdentityVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SBTIdentityVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.contract.Transact(opts, method, params...)
}

// SBTIDENTITYPROOFQUERYID is a free data retrieval call binding the contract method 0x2dff9de4.
//
// Solidity: function SBT_IDENTITY_PROOF_QUERY_ID() view returns(string)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) SBTIDENTITYPROOFQUERYID(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "SBT_IDENTITY_PROOF_QUERY_ID")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// SBTIDENTITYPROOFQUERYID is a free data retrieval call binding the contract method 0x2dff9de4.
//
// Solidity: function SBT_IDENTITY_PROOF_QUERY_ID() view returns(string)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) SBTIDENTITYPROOFQUERYID() (string, error) {
	return _SBTIdentityVerifier.Contract.SBTIDENTITYPROOFQUERYID(&_SBTIdentityVerifier.CallOpts)
}

// SBTIDENTITYPROOFQUERYID is a free data retrieval call binding the contract method 0x2dff9de4.
//
// Solidity: function SBT_IDENTITY_PROOF_QUERY_ID() view returns(string)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) SBTIDENTITYPROOFQUERYID() (string, error) {
	return _SBTIdentityVerifier.Contract.SBTIDENTITYPROOFQUERYID(&_SBTIdentityVerifier.CallOpts)
}

// AddressToIdentityId is a free data retrieval call binding the contract method 0xb4528ff9.
//
// Solidity: function addressToIdentityId(address ) view returns(uint256)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) AddressToIdentityId(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "addressToIdentityId", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AddressToIdentityId is a free data retrieval call binding the contract method 0xb4528ff9.
//
// Solidity: function addressToIdentityId(address ) view returns(uint256)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) AddressToIdentityId(arg0 common.Address) (*big.Int, error) {
	return _SBTIdentityVerifier.Contract.AddressToIdentityId(&_SBTIdentityVerifier.CallOpts, arg0)
}

// AddressToIdentityId is a free data retrieval call binding the contract method 0xb4528ff9.
//
// Solidity: function addressToIdentityId(address ) view returns(uint256)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) AddressToIdentityId(arg0 common.Address) (*big.Int, error) {
	return _SBTIdentityVerifier.Contract.AddressToIdentityId(&_SBTIdentityVerifier.CallOpts, arg0)
}

// GetAllowedIssuers is a free data retrieval call binding the contract method 0xef8dbdd3.
//
// Solidity: function getAllowedIssuers(uint256 schema_) view returns(uint256[])
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) GetAllowedIssuers(opts *bind.CallOpts, schema_ *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "getAllowedIssuers", schema_)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAllowedIssuers is a free data retrieval call binding the contract method 0xef8dbdd3.
//
// Solidity: function getAllowedIssuers(uint256 schema_) view returns(uint256[])
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) GetAllowedIssuers(schema_ *big.Int) ([]*big.Int, error) {
	return _SBTIdentityVerifier.Contract.GetAllowedIssuers(&_SBTIdentityVerifier.CallOpts, schema_)
}

// GetAllowedIssuers is a free data retrieval call binding the contract method 0xef8dbdd3.
//
// Solidity: function getAllowedIssuers(uint256 schema_) view returns(uint256[])
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) GetAllowedIssuers(schema_ *big.Int) ([]*big.Int, error) {
	return _SBTIdentityVerifier.Contract.GetAllowedIssuers(&_SBTIdentityVerifier.CallOpts, schema_)
}

// GetIdentityProofInfo is a free data retrieval call binding the contract method 0x5332d5ec.
//
// Solidity: function getIdentityProofInfo(uint256 identityId_) view returns((address,uint256,bool))
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) GetIdentityProofInfo(opts *bind.CallOpts, identityId_ *big.Int) (ISBTIdentityVerifierSBTIdentityProofInfo, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "getIdentityProofInfo", identityId_)

	if err != nil {
		return *new(ISBTIdentityVerifierSBTIdentityProofInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ISBTIdentityVerifierSBTIdentityProofInfo)).(*ISBTIdentityVerifierSBTIdentityProofInfo)

	return out0, err

}

// GetIdentityProofInfo is a free data retrieval call binding the contract method 0x5332d5ec.
//
// Solidity: function getIdentityProofInfo(uint256 identityId_) view returns((address,uint256,bool))
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) GetIdentityProofInfo(identityId_ *big.Int) (ISBTIdentityVerifierSBTIdentityProofInfo, error) {
	return _SBTIdentityVerifier.Contract.GetIdentityProofInfo(&_SBTIdentityVerifier.CallOpts, identityId_)
}

// GetIdentityProofInfo is a free data retrieval call binding the contract method 0x5332d5ec.
//
// Solidity: function getIdentityProofInfo(uint256 identityId_) view returns((address,uint256,bool))
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) GetIdentityProofInfo(identityId_ *big.Int) (ISBTIdentityVerifierSBTIdentityProofInfo, error) {
	return _SBTIdentityVerifier.Contract.GetIdentityProofInfo(&_SBTIdentityVerifier.CallOpts, identityId_)
}

// IsAllowedIssuer is a free data retrieval call binding the contract method 0x969f407e.
//
// Solidity: function isAllowedIssuer(uint256 schema_, uint256 issuerId_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) IsAllowedIssuer(opts *bind.CallOpts, schema_ *big.Int, issuerId_ *big.Int) (bool, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "isAllowedIssuer", schema_, issuerId_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAllowedIssuer is a free data retrieval call binding the contract method 0x969f407e.
//
// Solidity: function isAllowedIssuer(uint256 schema_, uint256 issuerId_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) IsAllowedIssuer(schema_ *big.Int, issuerId_ *big.Int) (bool, error) {
	return _SBTIdentityVerifier.Contract.IsAllowedIssuer(&_SBTIdentityVerifier.CallOpts, schema_, issuerId_)
}

// IsAllowedIssuer is a free data retrieval call binding the contract method 0x969f407e.
//
// Solidity: function isAllowedIssuer(uint256 schema_, uint256 issuerId_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) IsAllowedIssuer(schema_ *big.Int, issuerId_ *big.Int) (bool, error) {
	return _SBTIdentityVerifier.Contract.IsAllowedIssuer(&_SBTIdentityVerifier.CallOpts, schema_, issuerId_)
}

// IsIdentityProved is a free data retrieval call binding the contract method 0x413a3b43.
//
// Solidity: function isIdentityProved(address userAddr_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) IsIdentityProved(opts *bind.CallOpts, userAddr_ common.Address) (bool, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "isIdentityProved", userAddr_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsIdentityProved is a free data retrieval call binding the contract method 0x413a3b43.
//
// Solidity: function isIdentityProved(address userAddr_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) IsIdentityProved(userAddr_ common.Address) (bool, error) {
	return _SBTIdentityVerifier.Contract.IsIdentityProved(&_SBTIdentityVerifier.CallOpts, userAddr_)
}

// IsIdentityProved is a free data retrieval call binding the contract method 0x413a3b43.
//
// Solidity: function isIdentityProved(address userAddr_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) IsIdentityProved(userAddr_ common.Address) (bool, error) {
	return _SBTIdentityVerifier.Contract.IsIdentityProved(&_SBTIdentityVerifier.CallOpts, userAddr_)
}

// IsIdentityProved0 is a free data retrieval call binding the contract method 0x5428764d.
//
// Solidity: function isIdentityProved(uint256 identityId_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) IsIdentityProved0(opts *bind.CallOpts, identityId_ *big.Int) (bool, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "isIdentityProved0", identityId_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsIdentityProved0 is a free data retrieval call binding the contract method 0x5428764d.
//
// Solidity: function isIdentityProved(uint256 identityId_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) IsIdentityProved0(identityId_ *big.Int) (bool, error) {
	return _SBTIdentityVerifier.Contract.IsIdentityProved0(&_SBTIdentityVerifier.CallOpts, identityId_)
}

// IsIdentityProved0 is a free data retrieval call binding the contract method 0x5428764d.
//
// Solidity: function isIdentityProved(uint256 identityId_) view returns(bool)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) IsIdentityProved0(identityId_ *big.Int) (bool, error) {
	return _SBTIdentityVerifier.Contract.IsIdentityProved0(&_SBTIdentityVerifier.CallOpts, identityId_)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) Owner() (common.Address, error) {
	return _SBTIdentityVerifier.Contract.Owner(&_SBTIdentityVerifier.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) Owner() (common.Address, error) {
	return _SBTIdentityVerifier.Contract.Owner(&_SBTIdentityVerifier.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) ProxiableUUID() ([32]byte, error) {
	return _SBTIdentityVerifier.Contract.ProxiableUUID(&_SBTIdentityVerifier.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) ProxiableUUID() ([32]byte, error) {
	return _SBTIdentityVerifier.Contract.ProxiableUUID(&_SBTIdentityVerifier.CallOpts)
}

// SbtToken is a free data retrieval call binding the contract method 0xe359d997.
//
// Solidity: function sbtToken() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) SbtToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "sbtToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SbtToken is a free data retrieval call binding the contract method 0xe359d997.
//
// Solidity: function sbtToken() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) SbtToken() (common.Address, error) {
	return _SBTIdentityVerifier.Contract.SbtToken(&_SBTIdentityVerifier.CallOpts)
}

// SbtToken is a free data retrieval call binding the contract method 0xe359d997.
//
// Solidity: function sbtToken() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) SbtToken() (common.Address, error) {
	return _SBTIdentityVerifier.Contract.SbtToken(&_SBTIdentityVerifier.CallOpts)
}

// ZkpQueriesStorage is a free data retrieval call binding the contract method 0xb4db08cc.
//
// Solidity: function zkpQueriesStorage() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierCaller) ZkpQueriesStorage(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SBTIdentityVerifier.contract.Call(opts, &out, "zkpQueriesStorage")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ZkpQueriesStorage is a free data retrieval call binding the contract method 0xb4db08cc.
//
// Solidity: function zkpQueriesStorage() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) ZkpQueriesStorage() (common.Address, error) {
	return _SBTIdentityVerifier.Contract.ZkpQueriesStorage(&_SBTIdentityVerifier.CallOpts)
}

// ZkpQueriesStorage is a free data retrieval call binding the contract method 0xb4db08cc.
//
// Solidity: function zkpQueriesStorage() view returns(address)
func (_SBTIdentityVerifier *SBTIdentityVerifierCallerSession) ZkpQueriesStorage() (common.Address, error) {
	return _SBTIdentityVerifier.Contract.ZkpQueriesStorage(&_SBTIdentityVerifier.CallOpts)
}

// SBTIdentityVerifierInit is a paid mutator transaction binding the contract method 0x935454d3.
//
// Solidity: function __SBTIdentityVerifier_init(address zkpQueriesStorage_, address sbtToken_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) SBTIdentityVerifierInit(opts *bind.TransactOpts, zkpQueriesStorage_ common.Address, sbtToken_ common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "__SBTIdentityVerifier_init", zkpQueriesStorage_, sbtToken_)
}

// SBTIdentityVerifierInit is a paid mutator transaction binding the contract method 0x935454d3.
//
// Solidity: function __SBTIdentityVerifier_init(address zkpQueriesStorage_, address sbtToken_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) SBTIdentityVerifierInit(zkpQueriesStorage_ common.Address, sbtToken_ common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.SBTIdentityVerifierInit(&_SBTIdentityVerifier.TransactOpts, zkpQueriesStorage_, sbtToken_)
}

// SBTIdentityVerifierInit is a paid mutator transaction binding the contract method 0x935454d3.
//
// Solidity: function __SBTIdentityVerifier_init(address zkpQueriesStorage_, address sbtToken_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) SBTIdentityVerifierInit(zkpQueriesStorage_ common.Address, sbtToken_ common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.SBTIdentityVerifierInit(&_SBTIdentityVerifier.TransactOpts, zkpQueriesStorage_, sbtToken_)
}

// ProveIdentity is a paid mutator transaction binding the contract method 0x008b9130.
//
// Solidity: function proveIdentity(((uint256,uint256,uint256,bytes32[]),uint256[],uint256[2],uint256[2][2],uint256[2]) proveIdentityParams_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) ProveIdentity(opts *bind.TransactOpts, proveIdentityParams_ IBaseVerifierProveIdentityParams) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "proveIdentity", proveIdentityParams_)
}

// ProveIdentity is a paid mutator transaction binding the contract method 0x008b9130.
//
// Solidity: function proveIdentity(((uint256,uint256,uint256,bytes32[]),uint256[],uint256[2],uint256[2][2],uint256[2]) proveIdentityParams_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) ProveIdentity(proveIdentityParams_ IBaseVerifierProveIdentityParams) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.ProveIdentity(&_SBTIdentityVerifier.TransactOpts, proveIdentityParams_)
}

// ProveIdentity is a paid mutator transaction binding the contract method 0x008b9130.
//
// Solidity: function proveIdentity(((uint256,uint256,uint256,bytes32[]),uint256[],uint256[2],uint256[2][2],uint256[2]) proveIdentityParams_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) ProveIdentity(proveIdentityParams_ IBaseVerifierProveIdentityParams) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.ProveIdentity(&_SBTIdentityVerifier.TransactOpts, proveIdentityParams_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) RenounceOwnership() (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.RenounceOwnership(&_SBTIdentityVerifier.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.RenounceOwnership(&_SBTIdentityVerifier.TransactOpts)
}

// SetZKPQueriesStorage is a paid mutator transaction binding the contract method 0xddebe5c0.
//
// Solidity: function setZKPQueriesStorage(address newZKPQueriesStorage_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) SetZKPQueriesStorage(opts *bind.TransactOpts, newZKPQueriesStorage_ common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "setZKPQueriesStorage", newZKPQueriesStorage_)
}

// SetZKPQueriesStorage is a paid mutator transaction binding the contract method 0xddebe5c0.
//
// Solidity: function setZKPQueriesStorage(address newZKPQueriesStorage_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) SetZKPQueriesStorage(newZKPQueriesStorage_ common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.SetZKPQueriesStorage(&_SBTIdentityVerifier.TransactOpts, newZKPQueriesStorage_)
}

// SetZKPQueriesStorage is a paid mutator transaction binding the contract method 0xddebe5c0.
//
// Solidity: function setZKPQueriesStorage(address newZKPQueriesStorage_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) SetZKPQueriesStorage(newZKPQueriesStorage_ common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.SetZKPQueriesStorage(&_SBTIdentityVerifier.TransactOpts, newZKPQueriesStorage_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.TransferOwnership(&_SBTIdentityVerifier.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.TransferOwnership(&_SBTIdentityVerifier.TransactOpts, newOwner)
}

// TransitStateAndProveIdentity is a paid mutator transaction binding the contract method 0xd2fbb694.
//
// Solidity: function transitStateAndProveIdentity(((uint256,uint256,uint256,bytes32[]),uint256[],uint256[2],uint256[2][2],uint256[2]) proveIdentityParams_, (bytes32,(uint256,uint256),bytes) transitStateParams_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) TransitStateAndProveIdentity(opts *bind.TransactOpts, proveIdentityParams_ IBaseVerifierProveIdentityParams, transitStateParams_ IBaseVerifierTransitStateParams) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "transitStateAndProveIdentity", proveIdentityParams_, transitStateParams_)
}

// TransitStateAndProveIdentity is a paid mutator transaction binding the contract method 0xd2fbb694.
//
// Solidity: function transitStateAndProveIdentity(((uint256,uint256,uint256,bytes32[]),uint256[],uint256[2],uint256[2][2],uint256[2]) proveIdentityParams_, (bytes32,(uint256,uint256),bytes) transitStateParams_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) TransitStateAndProveIdentity(proveIdentityParams_ IBaseVerifierProveIdentityParams, transitStateParams_ IBaseVerifierTransitStateParams) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.TransitStateAndProveIdentity(&_SBTIdentityVerifier.TransactOpts, proveIdentityParams_, transitStateParams_)
}

// TransitStateAndProveIdentity is a paid mutator transaction binding the contract method 0xd2fbb694.
//
// Solidity: function transitStateAndProveIdentity(((uint256,uint256,uint256,bytes32[]),uint256[],uint256[2],uint256[2][2],uint256[2]) proveIdentityParams_, (bytes32,(uint256,uint256),bytes) transitStateParams_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) TransitStateAndProveIdentity(proveIdentityParams_ IBaseVerifierProveIdentityParams, transitStateParams_ IBaseVerifierTransitStateParams) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.TransitStateAndProveIdentity(&_SBTIdentityVerifier.TransactOpts, proveIdentityParams_, transitStateParams_)
}

// UpdateAllowedIssuers is a paid mutator transaction binding the contract method 0x051788a5.
//
// Solidity: function updateAllowedIssuers(uint256 schema_, uint256[] issuerIds_, bool isAdding_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) UpdateAllowedIssuers(opts *bind.TransactOpts, schema_ *big.Int, issuerIds_ []*big.Int, isAdding_ bool) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "updateAllowedIssuers", schema_, issuerIds_, isAdding_)
}

// UpdateAllowedIssuers is a paid mutator transaction binding the contract method 0x051788a5.
//
// Solidity: function updateAllowedIssuers(uint256 schema_, uint256[] issuerIds_, bool isAdding_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) UpdateAllowedIssuers(schema_ *big.Int, issuerIds_ []*big.Int, isAdding_ bool) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.UpdateAllowedIssuers(&_SBTIdentityVerifier.TransactOpts, schema_, issuerIds_, isAdding_)
}

// UpdateAllowedIssuers is a paid mutator transaction binding the contract method 0x051788a5.
//
// Solidity: function updateAllowedIssuers(uint256 schema_, uint256[] issuerIds_, bool isAdding_) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) UpdateAllowedIssuers(schema_ *big.Int, issuerIds_ []*big.Int, isAdding_ bool) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.UpdateAllowedIssuers(&_SBTIdentityVerifier.TransactOpts, schema_, issuerIds_, isAdding_)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.UpgradeTo(&_SBTIdentityVerifier.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.UpgradeTo(&_SBTIdentityVerifier.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SBTIdentityVerifier.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.UpgradeToAndCall(&_SBTIdentityVerifier.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SBTIdentityVerifier *SBTIdentityVerifierTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SBTIdentityVerifier.Contract.UpgradeToAndCall(&_SBTIdentityVerifier.TransactOpts, newImplementation, data)
}

// SBTIdentityVerifierAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierAdminChangedIterator struct {
	Event *SBTIdentityVerifierAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SBTIdentityVerifierAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SBTIdentityVerifierAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SBTIdentityVerifierAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SBTIdentityVerifierAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SBTIdentityVerifierAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SBTIdentityVerifierAdminChanged represents a AdminChanged event raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*SBTIdentityVerifierAdminChangedIterator, error) {

	logs, sub, err := _SBTIdentityVerifier.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierAdminChangedIterator{contract: _SBTIdentityVerifier.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *SBTIdentityVerifierAdminChanged) (event.Subscription, error) {

	logs, sub, err := _SBTIdentityVerifier.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SBTIdentityVerifierAdminChanged)
				if err := _SBTIdentityVerifier.contract.UnpackLog(event, "AdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) ParseAdminChanged(log types.Log) (*SBTIdentityVerifierAdminChanged, error) {
	event := new(SBTIdentityVerifierAdminChanged)
	if err := _SBTIdentityVerifier.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SBTIdentityVerifierBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierBeaconUpgradedIterator struct {
	Event *SBTIdentityVerifierBeaconUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SBTIdentityVerifierBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SBTIdentityVerifierBeaconUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SBTIdentityVerifierBeaconUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SBTIdentityVerifierBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SBTIdentityVerifierBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SBTIdentityVerifierBeaconUpgraded represents a BeaconUpgraded event raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*SBTIdentityVerifierBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierBeaconUpgradedIterator{contract: _SBTIdentityVerifier.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *SBTIdentityVerifierBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SBTIdentityVerifierBeaconUpgraded)
				if err := _SBTIdentityVerifier.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) ParseBeaconUpgraded(log types.Log) (*SBTIdentityVerifierBeaconUpgraded, error) {
	event := new(SBTIdentityVerifierBeaconUpgraded)
	if err := _SBTIdentityVerifier.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SBTIdentityVerifierInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierInitializedIterator struct {
	Event *SBTIdentityVerifierInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SBTIdentityVerifierInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SBTIdentityVerifierInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SBTIdentityVerifierInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SBTIdentityVerifierInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SBTIdentityVerifierInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SBTIdentityVerifierInitialized represents a Initialized event raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) FilterInitialized(opts *bind.FilterOpts) (*SBTIdentityVerifierInitializedIterator, error) {

	logs, sub, err := _SBTIdentityVerifier.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierInitializedIterator{contract: _SBTIdentityVerifier.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SBTIdentityVerifierInitialized) (event.Subscription, error) {

	logs, sub, err := _SBTIdentityVerifier.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SBTIdentityVerifierInitialized)
				if err := _SBTIdentityVerifier.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) ParseInitialized(log types.Log) (*SBTIdentityVerifierInitialized, error) {
	event := new(SBTIdentityVerifierInitialized)
	if err := _SBTIdentityVerifier.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SBTIdentityVerifierOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierOwnershipTransferredIterator struct {
	Event *SBTIdentityVerifierOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SBTIdentityVerifierOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SBTIdentityVerifierOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SBTIdentityVerifierOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SBTIdentityVerifierOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SBTIdentityVerifierOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SBTIdentityVerifierOwnershipTransferred represents a OwnershipTransferred event raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SBTIdentityVerifierOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierOwnershipTransferredIterator{contract: _SBTIdentityVerifier.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SBTIdentityVerifierOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SBTIdentityVerifierOwnershipTransferred)
				if err := _SBTIdentityVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) ParseOwnershipTransferred(log types.Log) (*SBTIdentityVerifierOwnershipTransferred, error) {
	event := new(SBTIdentityVerifierOwnershipTransferred)
	if err := _SBTIdentityVerifier.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SBTIdentityVerifierSBTIdentityProvedIterator is returned from FilterSBTIdentityProved and is used to iterate over the raw logs and unpacked data for SBTIdentityProved events raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierSBTIdentityProvedIterator struct {
	Event *SBTIdentityVerifierSBTIdentityProved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SBTIdentityVerifierSBTIdentityProvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SBTIdentityVerifierSBTIdentityProved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SBTIdentityVerifierSBTIdentityProved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SBTIdentityVerifierSBTIdentityProvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SBTIdentityVerifierSBTIdentityProvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SBTIdentityVerifierSBTIdentityProved represents a SBTIdentityProved event raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierSBTIdentityProved struct {
	IdentityId *big.Int
	SenderAddr common.Address
	TokenAddr  common.Address
	TokenID    *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSBTIdentityProved is a free log retrieval operation binding the contract event 0x33d963db6181e83a89d884ce98977f9cf447ffef066289de3ced2d2006441391.
//
// Solidity: event SBTIdentityProved(uint256 indexed identityId, address senderAddr, address tokenAddr, uint256 tokenID)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) FilterSBTIdentityProved(opts *bind.FilterOpts, identityId []*big.Int) (*SBTIdentityVerifierSBTIdentityProvedIterator, error) {

	var identityIdRule []interface{}
	for _, identityIdItem := range identityId {
		identityIdRule = append(identityIdRule, identityIdItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.FilterLogs(opts, "SBTIdentityProved", identityIdRule)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierSBTIdentityProvedIterator{contract: _SBTIdentityVerifier.contract, event: "SBTIdentityProved", logs: logs, sub: sub}, nil
}

// WatchSBTIdentityProved is a free log subscription operation binding the contract event 0x33d963db6181e83a89d884ce98977f9cf447ffef066289de3ced2d2006441391.
//
// Solidity: event SBTIdentityProved(uint256 indexed identityId, address senderAddr, address tokenAddr, uint256 tokenID)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) WatchSBTIdentityProved(opts *bind.WatchOpts, sink chan<- *SBTIdentityVerifierSBTIdentityProved, identityId []*big.Int) (event.Subscription, error) {

	var identityIdRule []interface{}
	for _, identityIdItem := range identityId {
		identityIdRule = append(identityIdRule, identityIdItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.WatchLogs(opts, "SBTIdentityProved", identityIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SBTIdentityVerifierSBTIdentityProved)
				if err := _SBTIdentityVerifier.contract.UnpackLog(event, "SBTIdentityProved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSBTIdentityProved is a log parse operation binding the contract event 0x33d963db6181e83a89d884ce98977f9cf447ffef066289de3ced2d2006441391.
//
// Solidity: event SBTIdentityProved(uint256 indexed identityId, address senderAddr, address tokenAddr, uint256 tokenID)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) ParseSBTIdentityProved(log types.Log) (*SBTIdentityVerifierSBTIdentityProved, error) {
	event := new(SBTIdentityVerifierSBTIdentityProved)
	if err := _SBTIdentityVerifier.contract.UnpackLog(event, "SBTIdentityProved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SBTIdentityVerifierUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierUpgradedIterator struct {
	Event *SBTIdentityVerifierUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SBTIdentityVerifierUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SBTIdentityVerifierUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SBTIdentityVerifierUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SBTIdentityVerifierUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SBTIdentityVerifierUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SBTIdentityVerifierUpgraded represents a Upgraded event raised by the SBTIdentityVerifier contract.
type SBTIdentityVerifierUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SBTIdentityVerifierUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SBTIdentityVerifierUpgradedIterator{contract: _SBTIdentityVerifier.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SBTIdentityVerifierUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SBTIdentityVerifier.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SBTIdentityVerifierUpgraded)
				if err := _SBTIdentityVerifier.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SBTIdentityVerifier *SBTIdentityVerifierFilterer) ParseUpgraded(log types.Log) (*SBTIdentityVerifierUpgraded, error) {
	event := new(SBTIdentityVerifierUpgraded)
	if err := _SBTIdentityVerifier.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
