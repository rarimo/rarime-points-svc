// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// Groth16VerifierHelperProofPoints is an auto generated low-level Go binding around an user-defined struct.
type Groth16VerifierHelperProofPoints struct {
	A [2]*big.Int
	B [2][2]*big.Int
	C [2]*big.Int
}

// SparseMerkleTreeNode is an auto generated low-level Go binding around an user-defined struct.
type SparseMerkleTreeNode struct {
	NodeType   uint8
	ChildLeft  uint64
	ChildRight uint64
	NodeHash   [32]byte
	Key        [32]byte
	Value      [32]byte
}

// SparseMerkleTreeProof is an auto generated low-level Go binding around an user-defined struct.
type SparseMerkleTreeProof struct {
	Root         [32]byte
	Siblings     [][32]byte
	Existence    bool
	Key          [32]byte
	Value        [32]byte
	AuxExistence bool
	AuxKey       [32]byte
	AuxValue     [32]byte
}

// FaceRegistryMetaData contains all meta data concerning the FaceRegistry contract.
var FaceRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedToCallVerifyProof\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"featureHash\",\"type\":\"uint256\"}],\"name\":\"FeatureHashAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"currentNonce\",\"type\":\"uint256\"}],\"name\":\"InvalidAccountNonce\",\"type\":\"error\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"internalType\":\"structGroth16VerifierHelper.ProofPoints\",\"name\":\"proof\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"pubSignals\",\"type\":\"uint256[]\"}],\"name\":\"InvalidCircomProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"KeyAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"maxDepth\",\"type\":\"uint32\"}],\"name\":\"MaxDepthExceedsHardCap\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxDepthIsZero\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxDepthReached\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"currentDepth\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"newDepth\",\"type\":\"uint32\"}],\"name\":\"NewMaxDepthMustBeLarger\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"NotAnOracle\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TreeAlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TreeIsNotEmpty\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TreeNotInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"UnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress\",\"type\":\"uint256\"}],\"name\":\"UserAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress\",\"type\":\"uint256\"}],\"name\":\"UserNotRegistered\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"MinThresholdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"newOwners\",\"type\":\"address[]\"}],\"name\":\"OwnersAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"removedOwners\",\"type\":\"address[]\"}],\"name\":\"OwnersRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"RootUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"userAddress\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newState\",\"type\":\"uint256\"}],\"name\":\"RulesUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldVerifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newVerifier\",\"type\":\"address\"}],\"name\":\"RulesVerifierUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"userAddress\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"featureHash\",\"type\":\"uint256\"}],\"name\":\"UserRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldVerifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newVerifier\",\"type\":\"address\"}],\"name\":\"VerifierUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"EVENT_ID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FACE_PROOF_SIGNALS_COUNT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ROOT_VALIDITY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RULES_PROOF_SIGNALS_COUNT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"evidenceRegistry_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"faceVerifier_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rulesVerifier_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minThreshold_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"treeHeight_\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"oracles_\",\"type\":\"address[]\"}],\"name\":\"__FaceRegistry_init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"oracles_\",\"type\":\"address[]\"}],\"name\":\"addOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"newOwners_\",\"type\":\"address[]\"}],\"name\":\"addOwners\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"evidenceRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"faceVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"}],\"name\":\"getFeatureHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"}],\"name\":\"getNodeByKey\",\"outputs\":[{\"components\":[{\"internalType\":\"enumSparseMerkleTree.NodeType\",\"name\":\"nodeType\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"childLeft\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"childRight\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"nodeHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"}],\"internalType\":\"structSparseMerkleTree.Node\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwners\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"}],\"name\":\"getProof\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"},{\"internalType\":\"bool\",\"name\":\"existence\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"value\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"auxExistence\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"auxKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"auxValue\",\"type\":\"bytes32\"}],\"internalType\":\"structSparseMerkleTree.Proof\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"}],\"name\":\"getRule\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"address_\",\"type\":\"uint256\"}],\"name\":\"getVerificationNonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"featureHash_\",\"type\":\"uint256\"}],\"name\":\"isFeatureHashUsed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle_\",\"type\":\"address\"}],\"name\":\"isOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"address_\",\"type\":\"address\"}],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root_\",\"type\":\"bytes32\"}],\"name\":\"isRootLatest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"root_\",\"type\":\"bytes32\"}],\"name\":\"isRootValid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"}],\"name\":\"isUserRegistered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"featureHash_\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"internalType\":\"structGroth16VerifierHelper.ProofPoints\",\"name\":\"zkPoints_\",\"type\":\"tuple\"}],\"name\":\"registerUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"oracles_\",\"type\":\"address[]\"}],\"name\":\"removeOracles\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"oldOwners_\",\"type\":\"address[]\"}],\"name\":\"removeOwners\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"rootHash\",\"type\":\"bytes32\"}],\"name\":\"roots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"featureHash\",\"type\":\"uint256\"}],\"name\":\"rules\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"state\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rulesVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newVerifier_\",\"type\":\"address\"}],\"name\":\"setFaceVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newThreshold_\",\"type\":\"uint256\"}],\"name\":\"setMinThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newVerifier_\",\"type\":\"address\"}],\"name\":\"setRulesVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newState_\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"internalType\":\"structGroth16VerifierHelper.ProofPoints\",\"name\":\"zkPoints_\",\"type\":\"tuple\"}],\"name\":\"updateRule\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"featureHash\",\"type\":\"uint256\"}],\"name\":\"usedFeatureHashes\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userAddress\",\"type\":\"uint256\"}],\"name\":\"userRegistryHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"featureHash\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// FaceRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use FaceRegistryMetaData.ABI instead.
var FaceRegistryABI = FaceRegistryMetaData.ABI

// FaceRegistry is an auto generated Go binding around an Ethereum contract.
type FaceRegistry struct {
	FaceRegistryCaller     // Read-only binding to the contract
	FaceRegistryTransactor // Write-only binding to the contract
	FaceRegistryFilterer   // Log filterer for contract events
}

// FaceRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type FaceRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FaceRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FaceRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FaceRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FaceRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FaceRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FaceRegistrySession struct {
	Contract     *FaceRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FaceRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FaceRegistryCallerSession struct {
	Contract *FaceRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// FaceRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FaceRegistryTransactorSession struct {
	Contract     *FaceRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// FaceRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type FaceRegistryRaw struct {
	Contract *FaceRegistry // Generic contract binding to access the raw methods on
}

// FaceRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FaceRegistryCallerRaw struct {
	Contract *FaceRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// FaceRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FaceRegistryTransactorRaw struct {
	Contract *FaceRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFaceRegistry creates a new instance of FaceRegistry, bound to a specific deployed contract.
func NewFaceRegistry(address common.Address, backend bind.ContractBackend) (*FaceRegistry, error) {
	contract, err := bindFaceRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FaceRegistry{FaceRegistryCaller: FaceRegistryCaller{contract: contract}, FaceRegistryTransactor: FaceRegistryTransactor{contract: contract}, FaceRegistryFilterer: FaceRegistryFilterer{contract: contract}}, nil
}

// NewFaceRegistryCaller creates a new read-only instance of FaceRegistry, bound to a specific deployed contract.
func NewFaceRegistryCaller(address common.Address, caller bind.ContractCaller) (*FaceRegistryCaller, error) {
	contract, err := bindFaceRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FaceRegistryCaller{contract: contract}, nil
}

// NewFaceRegistryTransactor creates a new write-only instance of FaceRegistry, bound to a specific deployed contract.
func NewFaceRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*FaceRegistryTransactor, error) {
	contract, err := bindFaceRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FaceRegistryTransactor{contract: contract}, nil
}

// NewFaceRegistryFilterer creates a new log filterer instance of FaceRegistry, bound to a specific deployed contract.
func NewFaceRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*FaceRegistryFilterer, error) {
	contract, err := bindFaceRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FaceRegistryFilterer{contract: contract}, nil
}

// bindFaceRegistry binds a generic wrapper to an already deployed contract.
func bindFaceRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FaceRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FaceRegistry *FaceRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FaceRegistry.Contract.FaceRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FaceRegistry *FaceRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FaceRegistry.Contract.FaceRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FaceRegistry *FaceRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FaceRegistry.Contract.FaceRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FaceRegistry *FaceRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FaceRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FaceRegistry *FaceRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FaceRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FaceRegistry *FaceRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FaceRegistry.Contract.contract.Transact(opts, method, params...)
}

// EVENTID is a free data retrieval call binding the contract method 0xb98efd4d.
//
// Solidity: function EVENT_ID() view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) EVENTID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "EVENT_ID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EVENTID is a free data retrieval call binding the contract method 0xb98efd4d.
//
// Solidity: function EVENT_ID() view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) EVENTID() (*big.Int, error) {
	return _FaceRegistry.Contract.EVENTID(&_FaceRegistry.CallOpts)
}

// EVENTID is a free data retrieval call binding the contract method 0xb98efd4d.
//
// Solidity: function EVENT_ID() view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) EVENTID() (*big.Int, error) {
	return _FaceRegistry.Contract.EVENTID(&_FaceRegistry.CallOpts)
}

// FACEPROOFSIGNALSCOUNT is a free data retrieval call binding the contract method 0xcd88e44b.
//
// Solidity: function FACE_PROOF_SIGNALS_COUNT() view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) FACEPROOFSIGNALSCOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "FACE_PROOF_SIGNALS_COUNT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FACEPROOFSIGNALSCOUNT is a free data retrieval call binding the contract method 0xcd88e44b.
//
// Solidity: function FACE_PROOF_SIGNALS_COUNT() view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) FACEPROOFSIGNALSCOUNT() (*big.Int, error) {
	return _FaceRegistry.Contract.FACEPROOFSIGNALSCOUNT(&_FaceRegistry.CallOpts)
}

// FACEPROOFSIGNALSCOUNT is a free data retrieval call binding the contract method 0xcd88e44b.
//
// Solidity: function FACE_PROOF_SIGNALS_COUNT() view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) FACEPROOFSIGNALSCOUNT() (*big.Int, error) {
	return _FaceRegistry.Contract.FACEPROOFSIGNALSCOUNT(&_FaceRegistry.CallOpts)
}

// ROOTVALIDITY is a free data retrieval call binding the contract method 0xcffe9676.
//
// Solidity: function ROOT_VALIDITY() view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) ROOTVALIDITY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "ROOT_VALIDITY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ROOTVALIDITY is a free data retrieval call binding the contract method 0xcffe9676.
//
// Solidity: function ROOT_VALIDITY() view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) ROOTVALIDITY() (*big.Int, error) {
	return _FaceRegistry.Contract.ROOTVALIDITY(&_FaceRegistry.CallOpts)
}

// ROOTVALIDITY is a free data retrieval call binding the contract method 0xcffe9676.
//
// Solidity: function ROOT_VALIDITY() view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) ROOTVALIDITY() (*big.Int, error) {
	return _FaceRegistry.Contract.ROOTVALIDITY(&_FaceRegistry.CallOpts)
}

// RULESPROOFSIGNALSCOUNT is a free data retrieval call binding the contract method 0xcb74c140.
//
// Solidity: function RULES_PROOF_SIGNALS_COUNT() view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) RULESPROOFSIGNALSCOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "RULES_PROOF_SIGNALS_COUNT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RULESPROOFSIGNALSCOUNT is a free data retrieval call binding the contract method 0xcb74c140.
//
// Solidity: function RULES_PROOF_SIGNALS_COUNT() view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) RULESPROOFSIGNALSCOUNT() (*big.Int, error) {
	return _FaceRegistry.Contract.RULESPROOFSIGNALSCOUNT(&_FaceRegistry.CallOpts)
}

// RULESPROOFSIGNALSCOUNT is a free data retrieval call binding the contract method 0xcb74c140.
//
// Solidity: function RULES_PROOF_SIGNALS_COUNT() view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) RULESPROOFSIGNALSCOUNT() (*big.Int, error) {
	return _FaceRegistry.Contract.RULESPROOFSIGNALSCOUNT(&_FaceRegistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FaceRegistry *FaceRegistryCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FaceRegistry *FaceRegistrySession) UPGRADEINTERFACEVERSION() (string, error) {
	return _FaceRegistry.Contract.UPGRADEINTERFACEVERSION(&_FaceRegistry.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FaceRegistry *FaceRegistryCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _FaceRegistry.Contract.UPGRADEINTERFACEVERSION(&_FaceRegistry.CallOpts)
}

// EvidenceRegistry is a free data retrieval call binding the contract method 0x95272e6d.
//
// Solidity: function evidenceRegistry() view returns(address)
func (_FaceRegistry *FaceRegistryCaller) EvidenceRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "evidenceRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EvidenceRegistry is a free data retrieval call binding the contract method 0x95272e6d.
//
// Solidity: function evidenceRegistry() view returns(address)
func (_FaceRegistry *FaceRegistrySession) EvidenceRegistry() (common.Address, error) {
	return _FaceRegistry.Contract.EvidenceRegistry(&_FaceRegistry.CallOpts)
}

// EvidenceRegistry is a free data retrieval call binding the contract method 0x95272e6d.
//
// Solidity: function evidenceRegistry() view returns(address)
func (_FaceRegistry *FaceRegistryCallerSession) EvidenceRegistry() (common.Address, error) {
	return _FaceRegistry.Contract.EvidenceRegistry(&_FaceRegistry.CallOpts)
}

// FaceVerifier is a free data retrieval call binding the contract method 0xe34d6431.
//
// Solidity: function faceVerifier() view returns(address)
func (_FaceRegistry *FaceRegistryCaller) FaceVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "faceVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FaceVerifier is a free data retrieval call binding the contract method 0xe34d6431.
//
// Solidity: function faceVerifier() view returns(address)
func (_FaceRegistry *FaceRegistrySession) FaceVerifier() (common.Address, error) {
	return _FaceRegistry.Contract.FaceVerifier(&_FaceRegistry.CallOpts)
}

// FaceVerifier is a free data retrieval call binding the contract method 0xe34d6431.
//
// Solidity: function faceVerifier() view returns(address)
func (_FaceRegistry *FaceRegistryCallerSession) FaceVerifier() (common.Address, error) {
	return _FaceRegistry.Contract.FaceVerifier(&_FaceRegistry.CallOpts)
}

// GetFeatureHash is a free data retrieval call binding the contract method 0x7bad5d93.
//
// Solidity: function getFeatureHash(uint256 userAddress_) view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) GetFeatureHash(opts *bind.CallOpts, userAddress_ *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getFeatureHash", userAddress_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeatureHash is a free data retrieval call binding the contract method 0x7bad5d93.
//
// Solidity: function getFeatureHash(uint256 userAddress_) view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) GetFeatureHash(userAddress_ *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.GetFeatureHash(&_FaceRegistry.CallOpts, userAddress_)
}

// GetFeatureHash is a free data retrieval call binding the contract method 0x7bad5d93.
//
// Solidity: function getFeatureHash(uint256 userAddress_) view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) GetFeatureHash(userAddress_ *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.GetFeatureHash(&_FaceRegistry.CallOpts, userAddress_)
}

// GetMinThreshold is a free data retrieval call binding the contract method 0xe6bbe9dd.
//
// Solidity: function getMinThreshold() view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) GetMinThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getMinThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinThreshold is a free data retrieval call binding the contract method 0xe6bbe9dd.
//
// Solidity: function getMinThreshold() view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) GetMinThreshold() (*big.Int, error) {
	return _FaceRegistry.Contract.GetMinThreshold(&_FaceRegistry.CallOpts)
}

// GetMinThreshold is a free data retrieval call binding the contract method 0xe6bbe9dd.
//
// Solidity: function getMinThreshold() view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) GetMinThreshold() (*big.Int, error) {
	return _FaceRegistry.Contract.GetMinThreshold(&_FaceRegistry.CallOpts)
}

// GetNodeByKey is a free data retrieval call binding the contract method 0x56cb5bb4.
//
// Solidity: function getNodeByKey(uint256 userAddress_) view returns((uint8,uint64,uint64,bytes32,bytes32,bytes32))
func (_FaceRegistry *FaceRegistryCaller) GetNodeByKey(opts *bind.CallOpts, userAddress_ *big.Int) (SparseMerkleTreeNode, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getNodeByKey", userAddress_)

	if err != nil {
		return *new(SparseMerkleTreeNode), err
	}

	out0 := *abi.ConvertType(out[0], new(SparseMerkleTreeNode)).(*SparseMerkleTreeNode)

	return out0, err

}

// GetNodeByKey is a free data retrieval call binding the contract method 0x56cb5bb4.
//
// Solidity: function getNodeByKey(uint256 userAddress_) view returns((uint8,uint64,uint64,bytes32,bytes32,bytes32))
func (_FaceRegistry *FaceRegistrySession) GetNodeByKey(userAddress_ *big.Int) (SparseMerkleTreeNode, error) {
	return _FaceRegistry.Contract.GetNodeByKey(&_FaceRegistry.CallOpts, userAddress_)
}

// GetNodeByKey is a free data retrieval call binding the contract method 0x56cb5bb4.
//
// Solidity: function getNodeByKey(uint256 userAddress_) view returns((uint8,uint64,uint64,bytes32,bytes32,bytes32))
func (_FaceRegistry *FaceRegistryCallerSession) GetNodeByKey(userAddress_ *big.Int) (SparseMerkleTreeNode, error) {
	return _FaceRegistry.Contract.GetNodeByKey(&_FaceRegistry.CallOpts, userAddress_)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_FaceRegistry *FaceRegistryCaller) GetOracles(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getOracles")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_FaceRegistry *FaceRegistrySession) GetOracles() ([]common.Address, error) {
	return _FaceRegistry.Contract.GetOracles(&_FaceRegistry.CallOpts)
}

// GetOracles is a free data retrieval call binding the contract method 0x40884c52.
//
// Solidity: function getOracles() view returns(address[])
func (_FaceRegistry *FaceRegistryCallerSession) GetOracles() ([]common.Address, error) {
	return _FaceRegistry.Contract.GetOracles(&_FaceRegistry.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_FaceRegistry *FaceRegistryCaller) GetOwners(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getOwners")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_FaceRegistry *FaceRegistrySession) GetOwners() ([]common.Address, error) {
	return _FaceRegistry.Contract.GetOwners(&_FaceRegistry.CallOpts)
}

// GetOwners is a free data retrieval call binding the contract method 0xa0e67e2b.
//
// Solidity: function getOwners() view returns(address[])
func (_FaceRegistry *FaceRegistryCallerSession) GetOwners() ([]common.Address, error) {
	return _FaceRegistry.Contract.GetOwners(&_FaceRegistry.CallOpts)
}

// GetProof is a free data retrieval call binding the contract method 0x11149ada.
//
// Solidity: function getProof(uint256 userAddress_) view returns((bytes32,bytes32[],bool,bytes32,bytes32,bool,bytes32,bytes32))
func (_FaceRegistry *FaceRegistryCaller) GetProof(opts *bind.CallOpts, userAddress_ *big.Int) (SparseMerkleTreeProof, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getProof", userAddress_)

	if err != nil {
		return *new(SparseMerkleTreeProof), err
	}

	out0 := *abi.ConvertType(out[0], new(SparseMerkleTreeProof)).(*SparseMerkleTreeProof)

	return out0, err

}

// GetProof is a free data retrieval call binding the contract method 0x11149ada.
//
// Solidity: function getProof(uint256 userAddress_) view returns((bytes32,bytes32[],bool,bytes32,bytes32,bool,bytes32,bytes32))
func (_FaceRegistry *FaceRegistrySession) GetProof(userAddress_ *big.Int) (SparseMerkleTreeProof, error) {
	return _FaceRegistry.Contract.GetProof(&_FaceRegistry.CallOpts, userAddress_)
}

// GetProof is a free data retrieval call binding the contract method 0x11149ada.
//
// Solidity: function getProof(uint256 userAddress_) view returns((bytes32,bytes32[],bool,bytes32,bytes32,bool,bytes32,bytes32))
func (_FaceRegistry *FaceRegistryCallerSession) GetProof(userAddress_ *big.Int) (SparseMerkleTreeProof, error) {
	return _FaceRegistry.Contract.GetProof(&_FaceRegistry.CallOpts, userAddress_)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_FaceRegistry *FaceRegistryCaller) GetRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_FaceRegistry *FaceRegistrySession) GetRoot() ([32]byte, error) {
	return _FaceRegistry.Contract.GetRoot(&_FaceRegistry.CallOpts)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_FaceRegistry *FaceRegistryCallerSession) GetRoot() ([32]byte, error) {
	return _FaceRegistry.Contract.GetRoot(&_FaceRegistry.CallOpts)
}

// GetRule is a free data retrieval call binding the contract method 0x48462037.
//
// Solidity: function getRule(uint256 userAddress_) view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) GetRule(opts *bind.CallOpts, userAddress_ *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getRule", userAddress_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRule is a free data retrieval call binding the contract method 0x48462037.
//
// Solidity: function getRule(uint256 userAddress_) view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) GetRule(userAddress_ *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.GetRule(&_FaceRegistry.CallOpts, userAddress_)
}

// GetRule is a free data retrieval call binding the contract method 0x48462037.
//
// Solidity: function getRule(uint256 userAddress_) view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) GetRule(userAddress_ *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.GetRule(&_FaceRegistry.CallOpts, userAddress_)
}

// GetVerificationNonce is a free data retrieval call binding the contract method 0x99bd4392.
//
// Solidity: function getVerificationNonce(uint256 address_) view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) GetVerificationNonce(opts *bind.CallOpts, address_ *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "getVerificationNonce", address_)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVerificationNonce is a free data retrieval call binding the contract method 0x99bd4392.
//
// Solidity: function getVerificationNonce(uint256 address_) view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) GetVerificationNonce(address_ *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.GetVerificationNonce(&_FaceRegistry.CallOpts, address_)
}

// GetVerificationNonce is a free data retrieval call binding the contract method 0x99bd4392.
//
// Solidity: function getVerificationNonce(uint256 address_) view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) GetVerificationNonce(address_ *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.GetVerificationNonce(&_FaceRegistry.CallOpts, address_)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_FaceRegistry *FaceRegistryCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_FaceRegistry *FaceRegistrySession) Implementation() (common.Address, error) {
	return _FaceRegistry.Contract.Implementation(&_FaceRegistry.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_FaceRegistry *FaceRegistryCallerSession) Implementation() (common.Address, error) {
	return _FaceRegistry.Contract.Implementation(&_FaceRegistry.CallOpts)
}

// IsFeatureHashUsed is a free data retrieval call binding the contract method 0xba7adecf.
//
// Solidity: function isFeatureHashUsed(uint256 featureHash_) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) IsFeatureHashUsed(opts *bind.CallOpts, featureHash_ *big.Int) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "isFeatureHashUsed", featureHash_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFeatureHashUsed is a free data retrieval call binding the contract method 0xba7adecf.
//
// Solidity: function isFeatureHashUsed(uint256 featureHash_) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) IsFeatureHashUsed(featureHash_ *big.Int) (bool, error) {
	return _FaceRegistry.Contract.IsFeatureHashUsed(&_FaceRegistry.CallOpts, featureHash_)
}

// IsFeatureHashUsed is a free data retrieval call binding the contract method 0xba7adecf.
//
// Solidity: function isFeatureHashUsed(uint256 featureHash_) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) IsFeatureHashUsed(featureHash_ *big.Int) (bool, error) {
	return _FaceRegistry.Contract.IsFeatureHashUsed(&_FaceRegistry.CallOpts, featureHash_)
}

// IsOracle is a free data retrieval call binding the contract method 0xa97e5c93.
//
// Solidity: function isOracle(address oracle_) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) IsOracle(opts *bind.CallOpts, oracle_ common.Address) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "isOracle", oracle_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracle is a free data retrieval call binding the contract method 0xa97e5c93.
//
// Solidity: function isOracle(address oracle_) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) IsOracle(oracle_ common.Address) (bool, error) {
	return _FaceRegistry.Contract.IsOracle(&_FaceRegistry.CallOpts, oracle_)
}

// IsOracle is a free data retrieval call binding the contract method 0xa97e5c93.
//
// Solidity: function isOracle(address oracle_) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) IsOracle(oracle_ common.Address) (bool, error) {
	return _FaceRegistry.Contract.IsOracle(&_FaceRegistry.CallOpts, oracle_)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address address_) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) IsOwner(opts *bind.CallOpts, address_ common.Address) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "isOwner", address_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address address_) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) IsOwner(address_ common.Address) (bool, error) {
	return _FaceRegistry.Contract.IsOwner(&_FaceRegistry.CallOpts, address_)
}

// IsOwner is a free data retrieval call binding the contract method 0x2f54bf6e.
//
// Solidity: function isOwner(address address_) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) IsOwner(address_ common.Address) (bool, error) {
	return _FaceRegistry.Contract.IsOwner(&_FaceRegistry.CallOpts, address_)
}

// IsRootLatest is a free data retrieval call binding the contract method 0x8492307f.
//
// Solidity: function isRootLatest(bytes32 root_) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) IsRootLatest(opts *bind.CallOpts, root_ [32]byte) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "isRootLatest", root_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRootLatest is a free data retrieval call binding the contract method 0x8492307f.
//
// Solidity: function isRootLatest(bytes32 root_) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) IsRootLatest(root_ [32]byte) (bool, error) {
	return _FaceRegistry.Contract.IsRootLatest(&_FaceRegistry.CallOpts, root_)
}

// IsRootLatest is a free data retrieval call binding the contract method 0x8492307f.
//
// Solidity: function isRootLatest(bytes32 root_) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) IsRootLatest(root_ [32]byte) (bool, error) {
	return _FaceRegistry.Contract.IsRootLatest(&_FaceRegistry.CallOpts, root_)
}

// IsRootValid is a free data retrieval call binding the contract method 0x30ef41b4.
//
// Solidity: function isRootValid(bytes32 root_) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) IsRootValid(opts *bind.CallOpts, root_ [32]byte) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "isRootValid", root_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsRootValid is a free data retrieval call binding the contract method 0x30ef41b4.
//
// Solidity: function isRootValid(bytes32 root_) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) IsRootValid(root_ [32]byte) (bool, error) {
	return _FaceRegistry.Contract.IsRootValid(&_FaceRegistry.CallOpts, root_)
}

// IsRootValid is a free data retrieval call binding the contract method 0x30ef41b4.
//
// Solidity: function isRootValid(bytes32 root_) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) IsRootValid(root_ [32]byte) (bool, error) {
	return _FaceRegistry.Contract.IsRootValid(&_FaceRegistry.CallOpts, root_)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0xe0a58ee1.
//
// Solidity: function isUserRegistered(uint256 userAddress_) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) IsUserRegistered(opts *bind.CallOpts, userAddress_ *big.Int) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "isUserRegistered", userAddress_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUserRegistered is a free data retrieval call binding the contract method 0xe0a58ee1.
//
// Solidity: function isUserRegistered(uint256 userAddress_) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) IsUserRegistered(userAddress_ *big.Int) (bool, error) {
	return _FaceRegistry.Contract.IsUserRegistered(&_FaceRegistry.CallOpts, userAddress_)
}

// IsUserRegistered is a free data retrieval call binding the contract method 0xe0a58ee1.
//
// Solidity: function isUserRegistered(uint256 userAddress_) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) IsUserRegistered(userAddress_ *big.Int) (bool, error) {
	return _FaceRegistry.Contract.IsUserRegistered(&_FaceRegistry.CallOpts, userAddress_)
}

// MinThreshold is a free data retrieval call binding the contract method 0xc85501bb.
//
// Solidity: function minThreshold() view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) MinThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "minThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinThreshold is a free data retrieval call binding the contract method 0xc85501bb.
//
// Solidity: function minThreshold() view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) MinThreshold() (*big.Int, error) {
	return _FaceRegistry.Contract.MinThreshold(&_FaceRegistry.CallOpts)
}

// MinThreshold is a free data retrieval call binding the contract method 0xc85501bb.
//
// Solidity: function minThreshold() view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) MinThreshold() (*big.Int, error) {
	return _FaceRegistry.Contract.MinThreshold(&_FaceRegistry.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_FaceRegistry *FaceRegistryCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "nonces", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_FaceRegistry *FaceRegistrySession) Nonces(owner common.Address) (*big.Int, error) {
	return _FaceRegistry.Contract.Nonces(&_FaceRegistry.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_FaceRegistry *FaceRegistryCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _FaceRegistry.Contract.Nonces(&_FaceRegistry.CallOpts, owner)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FaceRegistry *FaceRegistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FaceRegistry *FaceRegistrySession) ProxiableUUID() ([32]byte, error) {
	return _FaceRegistry.Contract.ProxiableUUID(&_FaceRegistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FaceRegistry *FaceRegistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _FaceRegistry.Contract.ProxiableUUID(&_FaceRegistry.CallOpts)
}

// Roots is a free data retrieval call binding the contract method 0xae6dead7.
//
// Solidity: function roots(bytes32 rootHash) view returns(uint256 timestamp)
func (_FaceRegistry *FaceRegistryCaller) Roots(opts *bind.CallOpts, rootHash [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "roots", rootHash)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Roots is a free data retrieval call binding the contract method 0xae6dead7.
//
// Solidity: function roots(bytes32 rootHash) view returns(uint256 timestamp)
func (_FaceRegistry *FaceRegistrySession) Roots(rootHash [32]byte) (*big.Int, error) {
	return _FaceRegistry.Contract.Roots(&_FaceRegistry.CallOpts, rootHash)
}

// Roots is a free data retrieval call binding the contract method 0xae6dead7.
//
// Solidity: function roots(bytes32 rootHash) view returns(uint256 timestamp)
func (_FaceRegistry *FaceRegistryCallerSession) Roots(rootHash [32]byte) (*big.Int, error) {
	return _FaceRegistry.Contract.Roots(&_FaceRegistry.CallOpts, rootHash)
}

// Rules is a free data retrieval call binding the contract method 0x04d6ded4.
//
// Solidity: function rules(uint256 featureHash) view returns(uint256 state)
func (_FaceRegistry *FaceRegistryCaller) Rules(opts *bind.CallOpts, featureHash *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "rules", featureHash)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rules is a free data retrieval call binding the contract method 0x04d6ded4.
//
// Solidity: function rules(uint256 featureHash) view returns(uint256 state)
func (_FaceRegistry *FaceRegistrySession) Rules(featureHash *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.Rules(&_FaceRegistry.CallOpts, featureHash)
}

// Rules is a free data retrieval call binding the contract method 0x04d6ded4.
//
// Solidity: function rules(uint256 featureHash) view returns(uint256 state)
func (_FaceRegistry *FaceRegistryCallerSession) Rules(featureHash *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.Rules(&_FaceRegistry.CallOpts, featureHash)
}

// RulesVerifier is a free data retrieval call binding the contract method 0x2fe62436.
//
// Solidity: function rulesVerifier() view returns(address)
func (_FaceRegistry *FaceRegistryCaller) RulesVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "rulesVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RulesVerifier is a free data retrieval call binding the contract method 0x2fe62436.
//
// Solidity: function rulesVerifier() view returns(address)
func (_FaceRegistry *FaceRegistrySession) RulesVerifier() (common.Address, error) {
	return _FaceRegistry.Contract.RulesVerifier(&_FaceRegistry.CallOpts)
}

// RulesVerifier is a free data retrieval call binding the contract method 0x2fe62436.
//
// Solidity: function rulesVerifier() view returns(address)
func (_FaceRegistry *FaceRegistryCallerSession) RulesVerifier() (common.Address, error) {
	return _FaceRegistry.Contract.RulesVerifier(&_FaceRegistry.CallOpts)
}

// UsedFeatureHashes is a free data retrieval call binding the contract method 0x45e99512.
//
// Solidity: function usedFeatureHashes(uint256 featureHash) view returns(bool)
func (_FaceRegistry *FaceRegistryCaller) UsedFeatureHashes(opts *bind.CallOpts, featureHash *big.Int) (bool, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "usedFeatureHashes", featureHash)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedFeatureHashes is a free data retrieval call binding the contract method 0x45e99512.
//
// Solidity: function usedFeatureHashes(uint256 featureHash) view returns(bool)
func (_FaceRegistry *FaceRegistrySession) UsedFeatureHashes(featureHash *big.Int) (bool, error) {
	return _FaceRegistry.Contract.UsedFeatureHashes(&_FaceRegistry.CallOpts, featureHash)
}

// UsedFeatureHashes is a free data retrieval call binding the contract method 0x45e99512.
//
// Solidity: function usedFeatureHashes(uint256 featureHash) view returns(bool)
func (_FaceRegistry *FaceRegistryCallerSession) UsedFeatureHashes(featureHash *big.Int) (bool, error) {
	return _FaceRegistry.Contract.UsedFeatureHashes(&_FaceRegistry.CallOpts, featureHash)
}

// UserRegistryHash is a free data retrieval call binding the contract method 0x226a6665.
//
// Solidity: function userRegistryHash(uint256 userAddress) view returns(uint256 featureHash)
func (_FaceRegistry *FaceRegistryCaller) UserRegistryHash(opts *bind.CallOpts, userAddress *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FaceRegistry.contract.Call(opts, &out, "userRegistryHash", userAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserRegistryHash is a free data retrieval call binding the contract method 0x226a6665.
//
// Solidity: function userRegistryHash(uint256 userAddress) view returns(uint256 featureHash)
func (_FaceRegistry *FaceRegistrySession) UserRegistryHash(userAddress *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.UserRegistryHash(&_FaceRegistry.CallOpts, userAddress)
}

// UserRegistryHash is a free data retrieval call binding the contract method 0x226a6665.
//
// Solidity: function userRegistryHash(uint256 userAddress) view returns(uint256 featureHash)
func (_FaceRegistry *FaceRegistryCallerSession) UserRegistryHash(userAddress *big.Int) (*big.Int, error) {
	return _FaceRegistry.Contract.UserRegistryHash(&_FaceRegistry.CallOpts, userAddress)
}

// FaceRegistryInit is a paid mutator transaction binding the contract method 0x0535bebb.
//
// Solidity: function __FaceRegistry_init(address evidenceRegistry_, address faceVerifier_, address rulesVerifier_, uint256 minThreshold_, uint256 treeHeight_, address[] oracles_) returns()
func (_FaceRegistry *FaceRegistryTransactor) FaceRegistryInit(opts *bind.TransactOpts, evidenceRegistry_ common.Address, faceVerifier_ common.Address, rulesVerifier_ common.Address, minThreshold_ *big.Int, treeHeight_ *big.Int, oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "__FaceRegistry_init", evidenceRegistry_, faceVerifier_, rulesVerifier_, minThreshold_, treeHeight_, oracles_)
}

// FaceRegistryInit is a paid mutator transaction binding the contract method 0x0535bebb.
//
// Solidity: function __FaceRegistry_init(address evidenceRegistry_, address faceVerifier_, address rulesVerifier_, uint256 minThreshold_, uint256 treeHeight_, address[] oracles_) returns()
func (_FaceRegistry *FaceRegistrySession) FaceRegistryInit(evidenceRegistry_ common.Address, faceVerifier_ common.Address, rulesVerifier_ common.Address, minThreshold_ *big.Int, treeHeight_ *big.Int, oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.FaceRegistryInit(&_FaceRegistry.TransactOpts, evidenceRegistry_, faceVerifier_, rulesVerifier_, minThreshold_, treeHeight_, oracles_)
}

// FaceRegistryInit is a paid mutator transaction binding the contract method 0x0535bebb.
//
// Solidity: function __FaceRegistry_init(address evidenceRegistry_, address faceVerifier_, address rulesVerifier_, uint256 minThreshold_, uint256 treeHeight_, address[] oracles_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) FaceRegistryInit(evidenceRegistry_ common.Address, faceVerifier_ common.Address, rulesVerifier_ common.Address, minThreshold_ *big.Int, treeHeight_ *big.Int, oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.FaceRegistryInit(&_FaceRegistry.TransactOpts, evidenceRegistry_, faceVerifier_, rulesVerifier_, minThreshold_, treeHeight_, oracles_)
}

// AddOracles is a paid mutator transaction binding the contract method 0x205b931e.
//
// Solidity: function addOracles(address[] oracles_) returns()
func (_FaceRegistry *FaceRegistryTransactor) AddOracles(opts *bind.TransactOpts, oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "addOracles", oracles_)
}

// AddOracles is a paid mutator transaction binding the contract method 0x205b931e.
//
// Solidity: function addOracles(address[] oracles_) returns()
func (_FaceRegistry *FaceRegistrySession) AddOracles(oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.AddOracles(&_FaceRegistry.TransactOpts, oracles_)
}

// AddOracles is a paid mutator transaction binding the contract method 0x205b931e.
//
// Solidity: function addOracles(address[] oracles_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) AddOracles(oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.AddOracles(&_FaceRegistry.TransactOpts, oracles_)
}

// AddOwners is a paid mutator transaction binding the contract method 0x6c46a2c5.
//
// Solidity: function addOwners(address[] newOwners_) returns()
func (_FaceRegistry *FaceRegistryTransactor) AddOwners(opts *bind.TransactOpts, newOwners_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "addOwners", newOwners_)
}

// AddOwners is a paid mutator transaction binding the contract method 0x6c46a2c5.
//
// Solidity: function addOwners(address[] newOwners_) returns()
func (_FaceRegistry *FaceRegistrySession) AddOwners(newOwners_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.AddOwners(&_FaceRegistry.TransactOpts, newOwners_)
}

// AddOwners is a paid mutator transaction binding the contract method 0x6c46a2c5.
//
// Solidity: function addOwners(address[] newOwners_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) AddOwners(newOwners_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.AddOwners(&_FaceRegistry.TransactOpts, newOwners_)
}

// RegisterUser is a paid mutator transaction binding the contract method 0xb5c718d7.
//
// Solidity: function registerUser(uint256 userAddress_, uint256 featureHash_, (uint256[2],uint256[2][2],uint256[2]) zkPoints_) returns()
func (_FaceRegistry *FaceRegistryTransactor) RegisterUser(opts *bind.TransactOpts, userAddress_ *big.Int, featureHash_ *big.Int, zkPoints_ Groth16VerifierHelperProofPoints) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "registerUser", userAddress_, featureHash_, zkPoints_)
}

// RegisterUser is a paid mutator transaction binding the contract method 0xb5c718d7.
//
// Solidity: function registerUser(uint256 userAddress_, uint256 featureHash_, (uint256[2],uint256[2][2],uint256[2]) zkPoints_) returns()
func (_FaceRegistry *FaceRegistrySession) RegisterUser(userAddress_ *big.Int, featureHash_ *big.Int, zkPoints_ Groth16VerifierHelperProofPoints) (*types.Transaction, error) {
	return _FaceRegistry.Contract.RegisterUser(&_FaceRegistry.TransactOpts, userAddress_, featureHash_, zkPoints_)
}

// RegisterUser is a paid mutator transaction binding the contract method 0xb5c718d7.
//
// Solidity: function registerUser(uint256 userAddress_, uint256 featureHash_, (uint256[2],uint256[2][2],uint256[2]) zkPoints_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) RegisterUser(userAddress_ *big.Int, featureHash_ *big.Int, zkPoints_ Groth16VerifierHelperProofPoints) (*types.Transaction, error) {
	return _FaceRegistry.Contract.RegisterUser(&_FaceRegistry.TransactOpts, userAddress_, featureHash_, zkPoints_)
}

// RemoveOracles is a paid mutator transaction binding the contract method 0x45644fd6.
//
// Solidity: function removeOracles(address[] oracles_) returns()
func (_FaceRegistry *FaceRegistryTransactor) RemoveOracles(opts *bind.TransactOpts, oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "removeOracles", oracles_)
}

// RemoveOracles is a paid mutator transaction binding the contract method 0x45644fd6.
//
// Solidity: function removeOracles(address[] oracles_) returns()
func (_FaceRegistry *FaceRegistrySession) RemoveOracles(oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.RemoveOracles(&_FaceRegistry.TransactOpts, oracles_)
}

// RemoveOracles is a paid mutator transaction binding the contract method 0x45644fd6.
//
// Solidity: function removeOracles(address[] oracles_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) RemoveOracles(oracles_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.RemoveOracles(&_FaceRegistry.TransactOpts, oracles_)
}

// RemoveOwners is a paid mutator transaction binding the contract method 0xa9a5e3af.
//
// Solidity: function removeOwners(address[] oldOwners_) returns()
func (_FaceRegistry *FaceRegistryTransactor) RemoveOwners(opts *bind.TransactOpts, oldOwners_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "removeOwners", oldOwners_)
}

// RemoveOwners is a paid mutator transaction binding the contract method 0xa9a5e3af.
//
// Solidity: function removeOwners(address[] oldOwners_) returns()
func (_FaceRegistry *FaceRegistrySession) RemoveOwners(oldOwners_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.RemoveOwners(&_FaceRegistry.TransactOpts, oldOwners_)
}

// RemoveOwners is a paid mutator transaction binding the contract method 0xa9a5e3af.
//
// Solidity: function removeOwners(address[] oldOwners_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) RemoveOwners(oldOwners_ []common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.RemoveOwners(&_FaceRegistry.TransactOpts, oldOwners_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FaceRegistry *FaceRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FaceRegistry *FaceRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _FaceRegistry.Contract.RenounceOwnership(&_FaceRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FaceRegistry *FaceRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FaceRegistry.Contract.RenounceOwnership(&_FaceRegistry.TransactOpts)
}

// SetFaceVerifier is a paid mutator transaction binding the contract method 0xd0129e46.
//
// Solidity: function setFaceVerifier(address newVerifier_) returns()
func (_FaceRegistry *FaceRegistryTransactor) SetFaceVerifier(opts *bind.TransactOpts, newVerifier_ common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "setFaceVerifier", newVerifier_)
}

// SetFaceVerifier is a paid mutator transaction binding the contract method 0xd0129e46.
//
// Solidity: function setFaceVerifier(address newVerifier_) returns()
func (_FaceRegistry *FaceRegistrySession) SetFaceVerifier(newVerifier_ common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.SetFaceVerifier(&_FaceRegistry.TransactOpts, newVerifier_)
}

// SetFaceVerifier is a paid mutator transaction binding the contract method 0xd0129e46.
//
// Solidity: function setFaceVerifier(address newVerifier_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) SetFaceVerifier(newVerifier_ common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.SetFaceVerifier(&_FaceRegistry.TransactOpts, newVerifier_)
}

// SetMinThreshold is a paid mutator transaction binding the contract method 0x7f39a939.
//
// Solidity: function setMinThreshold(uint256 newThreshold_) returns()
func (_FaceRegistry *FaceRegistryTransactor) SetMinThreshold(opts *bind.TransactOpts, newThreshold_ *big.Int) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "setMinThreshold", newThreshold_)
}

// SetMinThreshold is a paid mutator transaction binding the contract method 0x7f39a939.
//
// Solidity: function setMinThreshold(uint256 newThreshold_) returns()
func (_FaceRegistry *FaceRegistrySession) SetMinThreshold(newThreshold_ *big.Int) (*types.Transaction, error) {
	return _FaceRegistry.Contract.SetMinThreshold(&_FaceRegistry.TransactOpts, newThreshold_)
}

// SetMinThreshold is a paid mutator transaction binding the contract method 0x7f39a939.
//
// Solidity: function setMinThreshold(uint256 newThreshold_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) SetMinThreshold(newThreshold_ *big.Int) (*types.Transaction, error) {
	return _FaceRegistry.Contract.SetMinThreshold(&_FaceRegistry.TransactOpts, newThreshold_)
}

// SetRulesVerifier is a paid mutator transaction binding the contract method 0xc319f413.
//
// Solidity: function setRulesVerifier(address newVerifier_) returns()
func (_FaceRegistry *FaceRegistryTransactor) SetRulesVerifier(opts *bind.TransactOpts, newVerifier_ common.Address) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "setRulesVerifier", newVerifier_)
}

// SetRulesVerifier is a paid mutator transaction binding the contract method 0xc319f413.
//
// Solidity: function setRulesVerifier(address newVerifier_) returns()
func (_FaceRegistry *FaceRegistrySession) SetRulesVerifier(newVerifier_ common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.SetRulesVerifier(&_FaceRegistry.TransactOpts, newVerifier_)
}

// SetRulesVerifier is a paid mutator transaction binding the contract method 0xc319f413.
//
// Solidity: function setRulesVerifier(address newVerifier_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) SetRulesVerifier(newVerifier_ common.Address) (*types.Transaction, error) {
	return _FaceRegistry.Contract.SetRulesVerifier(&_FaceRegistry.TransactOpts, newVerifier_)
}

// UpdateRule is a paid mutator transaction binding the contract method 0x3e62651b.
//
// Solidity: function updateRule(uint256 userAddress_, uint256 newState_, (uint256[2],uint256[2][2],uint256[2]) zkPoints_) returns()
func (_FaceRegistry *FaceRegistryTransactor) UpdateRule(opts *bind.TransactOpts, userAddress_ *big.Int, newState_ *big.Int, zkPoints_ Groth16VerifierHelperProofPoints) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "updateRule", userAddress_, newState_, zkPoints_)
}

// UpdateRule is a paid mutator transaction binding the contract method 0x3e62651b.
//
// Solidity: function updateRule(uint256 userAddress_, uint256 newState_, (uint256[2],uint256[2][2],uint256[2]) zkPoints_) returns()
func (_FaceRegistry *FaceRegistrySession) UpdateRule(userAddress_ *big.Int, newState_ *big.Int, zkPoints_ Groth16VerifierHelperProofPoints) (*types.Transaction, error) {
	return _FaceRegistry.Contract.UpdateRule(&_FaceRegistry.TransactOpts, userAddress_, newState_, zkPoints_)
}

// UpdateRule is a paid mutator transaction binding the contract method 0x3e62651b.
//
// Solidity: function updateRule(uint256 userAddress_, uint256 newState_, (uint256[2],uint256[2][2],uint256[2]) zkPoints_) returns()
func (_FaceRegistry *FaceRegistryTransactorSession) UpdateRule(userAddress_ *big.Int, newState_ *big.Int, zkPoints_ Groth16VerifierHelperProofPoints) (*types.Transaction, error) {
	return _FaceRegistry.Contract.UpdateRule(&_FaceRegistry.TransactOpts, userAddress_, newState_, zkPoints_)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FaceRegistry *FaceRegistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FaceRegistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FaceRegistry *FaceRegistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FaceRegistry.Contract.UpgradeToAndCall(&_FaceRegistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FaceRegistry *FaceRegistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FaceRegistry.Contract.UpgradeToAndCall(&_FaceRegistry.TransactOpts, newImplementation, data)
}

// FaceRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FaceRegistry contract.
type FaceRegistryInitializedIterator struct {
	Event *FaceRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *FaceRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryInitialized)
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
		it.Event = new(FaceRegistryInitialized)
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
func (it *FaceRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryInitialized represents a Initialized event raised by the FaceRegistry contract.
type FaceRegistryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FaceRegistry *FaceRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*FaceRegistryInitializedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryInitializedIterator{contract: _FaceRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FaceRegistry *FaceRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FaceRegistryInitialized) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryInitialized)
				if err := _FaceRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FaceRegistry *FaceRegistryFilterer) ParseInitialized(log types.Log) (*FaceRegistryInitialized, error) {
	event := new(FaceRegistryInitialized)
	if err := _FaceRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryMinThresholdUpdatedIterator is returned from FilterMinThresholdUpdated and is used to iterate over the raw logs and unpacked data for MinThresholdUpdated events raised by the FaceRegistry contract.
type FaceRegistryMinThresholdUpdatedIterator struct {
	Event *FaceRegistryMinThresholdUpdated // Event containing the contract specifics and raw log

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
func (it *FaceRegistryMinThresholdUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryMinThresholdUpdated)
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
		it.Event = new(FaceRegistryMinThresholdUpdated)
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
func (it *FaceRegistryMinThresholdUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryMinThresholdUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryMinThresholdUpdated represents a MinThresholdUpdated event raised by the FaceRegistry contract.
type FaceRegistryMinThresholdUpdated struct {
	OldThreshold *big.Int
	NewThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinThresholdUpdated is a free log retrieval operation binding the contract event 0xefca4185725c2f876881025e43317d8a3023f9f3c1b8926941d399ca79ae2eb5.
//
// Solidity: event MinThresholdUpdated(uint256 oldThreshold, uint256 newThreshold)
func (_FaceRegistry *FaceRegistryFilterer) FilterMinThresholdUpdated(opts *bind.FilterOpts) (*FaceRegistryMinThresholdUpdatedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "MinThresholdUpdated")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryMinThresholdUpdatedIterator{contract: _FaceRegistry.contract, event: "MinThresholdUpdated", logs: logs, sub: sub}, nil
}

// WatchMinThresholdUpdated is a free log subscription operation binding the contract event 0xefca4185725c2f876881025e43317d8a3023f9f3c1b8926941d399ca79ae2eb5.
//
// Solidity: event MinThresholdUpdated(uint256 oldThreshold, uint256 newThreshold)
func (_FaceRegistry *FaceRegistryFilterer) WatchMinThresholdUpdated(opts *bind.WatchOpts, sink chan<- *FaceRegistryMinThresholdUpdated) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "MinThresholdUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryMinThresholdUpdated)
				if err := _FaceRegistry.contract.UnpackLog(event, "MinThresholdUpdated", log); err != nil {
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

// ParseMinThresholdUpdated is a log parse operation binding the contract event 0xefca4185725c2f876881025e43317d8a3023f9f3c1b8926941d399ca79ae2eb5.
//
// Solidity: event MinThresholdUpdated(uint256 oldThreshold, uint256 newThreshold)
func (_FaceRegistry *FaceRegistryFilterer) ParseMinThresholdUpdated(log types.Log) (*FaceRegistryMinThresholdUpdated, error) {
	event := new(FaceRegistryMinThresholdUpdated)
	if err := _FaceRegistry.contract.UnpackLog(event, "MinThresholdUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryOwnersAddedIterator is returned from FilterOwnersAdded and is used to iterate over the raw logs and unpacked data for OwnersAdded events raised by the FaceRegistry contract.
type FaceRegistryOwnersAddedIterator struct {
	Event *FaceRegistryOwnersAdded // Event containing the contract specifics and raw log

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
func (it *FaceRegistryOwnersAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryOwnersAdded)
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
		it.Event = new(FaceRegistryOwnersAdded)
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
func (it *FaceRegistryOwnersAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryOwnersAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryOwnersAdded represents a OwnersAdded event raised by the FaceRegistry contract.
type FaceRegistryOwnersAdded struct {
	NewOwners []common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOwnersAdded is a free log retrieval operation binding the contract event 0x5fd1e185ef572e7f662fcc63b7c9e778b996190372868af5fe137132c811398e.
//
// Solidity: event OwnersAdded(address[] newOwners)
func (_FaceRegistry *FaceRegistryFilterer) FilterOwnersAdded(opts *bind.FilterOpts) (*FaceRegistryOwnersAddedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "OwnersAdded")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryOwnersAddedIterator{contract: _FaceRegistry.contract, event: "OwnersAdded", logs: logs, sub: sub}, nil
}

// WatchOwnersAdded is a free log subscription operation binding the contract event 0x5fd1e185ef572e7f662fcc63b7c9e778b996190372868af5fe137132c811398e.
//
// Solidity: event OwnersAdded(address[] newOwners)
func (_FaceRegistry *FaceRegistryFilterer) WatchOwnersAdded(opts *bind.WatchOpts, sink chan<- *FaceRegistryOwnersAdded) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "OwnersAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryOwnersAdded)
				if err := _FaceRegistry.contract.UnpackLog(event, "OwnersAdded", log); err != nil {
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

// ParseOwnersAdded is a log parse operation binding the contract event 0x5fd1e185ef572e7f662fcc63b7c9e778b996190372868af5fe137132c811398e.
//
// Solidity: event OwnersAdded(address[] newOwners)
func (_FaceRegistry *FaceRegistryFilterer) ParseOwnersAdded(log types.Log) (*FaceRegistryOwnersAdded, error) {
	event := new(FaceRegistryOwnersAdded)
	if err := _FaceRegistry.contract.UnpackLog(event, "OwnersAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryOwnersRemovedIterator is returned from FilterOwnersRemoved and is used to iterate over the raw logs and unpacked data for OwnersRemoved events raised by the FaceRegistry contract.
type FaceRegistryOwnersRemovedIterator struct {
	Event *FaceRegistryOwnersRemoved // Event containing the contract specifics and raw log

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
func (it *FaceRegistryOwnersRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryOwnersRemoved)
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
		it.Event = new(FaceRegistryOwnersRemoved)
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
func (it *FaceRegistryOwnersRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryOwnersRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryOwnersRemoved represents a OwnersRemoved event raised by the FaceRegistry contract.
type FaceRegistryOwnersRemoved struct {
	RemovedOwners []common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnersRemoved is a free log retrieval operation binding the contract event 0x0bbb8c3531454b5141cebfe14eba43275a099c31e3357a4653412a08b05ce0cc.
//
// Solidity: event OwnersRemoved(address[] removedOwners)
func (_FaceRegistry *FaceRegistryFilterer) FilterOwnersRemoved(opts *bind.FilterOpts) (*FaceRegistryOwnersRemovedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "OwnersRemoved")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryOwnersRemovedIterator{contract: _FaceRegistry.contract, event: "OwnersRemoved", logs: logs, sub: sub}, nil
}

// WatchOwnersRemoved is a free log subscription operation binding the contract event 0x0bbb8c3531454b5141cebfe14eba43275a099c31e3357a4653412a08b05ce0cc.
//
// Solidity: event OwnersRemoved(address[] removedOwners)
func (_FaceRegistry *FaceRegistryFilterer) WatchOwnersRemoved(opts *bind.WatchOpts, sink chan<- *FaceRegistryOwnersRemoved) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "OwnersRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryOwnersRemoved)
				if err := _FaceRegistry.contract.UnpackLog(event, "OwnersRemoved", log); err != nil {
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

// ParseOwnersRemoved is a log parse operation binding the contract event 0x0bbb8c3531454b5141cebfe14eba43275a099c31e3357a4653412a08b05ce0cc.
//
// Solidity: event OwnersRemoved(address[] removedOwners)
func (_FaceRegistry *FaceRegistryFilterer) ParseOwnersRemoved(log types.Log) (*FaceRegistryOwnersRemoved, error) {
	event := new(FaceRegistryOwnersRemoved)
	if err := _FaceRegistry.contract.UnpackLog(event, "OwnersRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryRootUpdatedIterator is returned from FilterRootUpdated and is used to iterate over the raw logs and unpacked data for RootUpdated events raised by the FaceRegistry contract.
type FaceRegistryRootUpdatedIterator struct {
	Event *FaceRegistryRootUpdated // Event containing the contract specifics and raw log

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
func (it *FaceRegistryRootUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryRootUpdated)
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
		it.Event = new(FaceRegistryRootUpdated)
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
func (it *FaceRegistryRootUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryRootUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryRootUpdated represents a RootUpdated event raised by the FaceRegistry contract.
type FaceRegistryRootUpdated struct {
	Root [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRootUpdated is a free log retrieval operation binding the contract event 0x2cbc14f49c068133583f7cb530018af451c87c1cf1327cf2a4ff4698c4730aa4.
//
// Solidity: event RootUpdated(bytes32 root)
func (_FaceRegistry *FaceRegistryFilterer) FilterRootUpdated(opts *bind.FilterOpts) (*FaceRegistryRootUpdatedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "RootUpdated")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryRootUpdatedIterator{contract: _FaceRegistry.contract, event: "RootUpdated", logs: logs, sub: sub}, nil
}

// WatchRootUpdated is a free log subscription operation binding the contract event 0x2cbc14f49c068133583f7cb530018af451c87c1cf1327cf2a4ff4698c4730aa4.
//
// Solidity: event RootUpdated(bytes32 root)
func (_FaceRegistry *FaceRegistryFilterer) WatchRootUpdated(opts *bind.WatchOpts, sink chan<- *FaceRegistryRootUpdated) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "RootUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryRootUpdated)
				if err := _FaceRegistry.contract.UnpackLog(event, "RootUpdated", log); err != nil {
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

// ParseRootUpdated is a log parse operation binding the contract event 0x2cbc14f49c068133583f7cb530018af451c87c1cf1327cf2a4ff4698c4730aa4.
//
// Solidity: event RootUpdated(bytes32 root)
func (_FaceRegistry *FaceRegistryFilterer) ParseRootUpdated(log types.Log) (*FaceRegistryRootUpdated, error) {
	event := new(FaceRegistryRootUpdated)
	if err := _FaceRegistry.contract.UnpackLog(event, "RootUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryRulesUpdatedIterator is returned from FilterRulesUpdated and is used to iterate over the raw logs and unpacked data for RulesUpdated events raised by the FaceRegistry contract.
type FaceRegistryRulesUpdatedIterator struct {
	Event *FaceRegistryRulesUpdated // Event containing the contract specifics and raw log

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
func (it *FaceRegistryRulesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryRulesUpdated)
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
		it.Event = new(FaceRegistryRulesUpdated)
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
func (it *FaceRegistryRulesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryRulesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryRulesUpdated represents a RulesUpdated event raised by the FaceRegistry contract.
type FaceRegistryRulesUpdated struct {
	UserAddress *big.Int
	NewState    *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRulesUpdated is a free log retrieval operation binding the contract event 0xdaacdc06082cd99bd0e888138e5010f2d105c79419e19d33482fe59b5572ddb1.
//
// Solidity: event RulesUpdated(uint256 userAddress, uint256 newState)
func (_FaceRegistry *FaceRegistryFilterer) FilterRulesUpdated(opts *bind.FilterOpts) (*FaceRegistryRulesUpdatedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "RulesUpdated")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryRulesUpdatedIterator{contract: _FaceRegistry.contract, event: "RulesUpdated", logs: logs, sub: sub}, nil
}

// WatchRulesUpdated is a free log subscription operation binding the contract event 0xdaacdc06082cd99bd0e888138e5010f2d105c79419e19d33482fe59b5572ddb1.
//
// Solidity: event RulesUpdated(uint256 userAddress, uint256 newState)
func (_FaceRegistry *FaceRegistryFilterer) WatchRulesUpdated(opts *bind.WatchOpts, sink chan<- *FaceRegistryRulesUpdated) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "RulesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryRulesUpdated)
				if err := _FaceRegistry.contract.UnpackLog(event, "RulesUpdated", log); err != nil {
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

// ParseRulesUpdated is a log parse operation binding the contract event 0xdaacdc06082cd99bd0e888138e5010f2d105c79419e19d33482fe59b5572ddb1.
//
// Solidity: event RulesUpdated(uint256 userAddress, uint256 newState)
func (_FaceRegistry *FaceRegistryFilterer) ParseRulesUpdated(log types.Log) (*FaceRegistryRulesUpdated, error) {
	event := new(FaceRegistryRulesUpdated)
	if err := _FaceRegistry.contract.UnpackLog(event, "RulesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryRulesVerifierUpdatedIterator is returned from FilterRulesVerifierUpdated and is used to iterate over the raw logs and unpacked data for RulesVerifierUpdated events raised by the FaceRegistry contract.
type FaceRegistryRulesVerifierUpdatedIterator struct {
	Event *FaceRegistryRulesVerifierUpdated // Event containing the contract specifics and raw log

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
func (it *FaceRegistryRulesVerifierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryRulesVerifierUpdated)
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
		it.Event = new(FaceRegistryRulesVerifierUpdated)
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
func (it *FaceRegistryRulesVerifierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryRulesVerifierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryRulesVerifierUpdated represents a RulesVerifierUpdated event raised by the FaceRegistry contract.
type FaceRegistryRulesVerifierUpdated struct {
	OldVerifier common.Address
	NewVerifier common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRulesVerifierUpdated is a free log retrieval operation binding the contract event 0x5a8b13378dd1cabe2f3617497f917ede72540af55149036f6afc92c031efbf1f.
//
// Solidity: event RulesVerifierUpdated(address oldVerifier, address newVerifier)
func (_FaceRegistry *FaceRegistryFilterer) FilterRulesVerifierUpdated(opts *bind.FilterOpts) (*FaceRegistryRulesVerifierUpdatedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "RulesVerifierUpdated")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryRulesVerifierUpdatedIterator{contract: _FaceRegistry.contract, event: "RulesVerifierUpdated", logs: logs, sub: sub}, nil
}

// WatchRulesVerifierUpdated is a free log subscription operation binding the contract event 0x5a8b13378dd1cabe2f3617497f917ede72540af55149036f6afc92c031efbf1f.
//
// Solidity: event RulesVerifierUpdated(address oldVerifier, address newVerifier)
func (_FaceRegistry *FaceRegistryFilterer) WatchRulesVerifierUpdated(opts *bind.WatchOpts, sink chan<- *FaceRegistryRulesVerifierUpdated) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "RulesVerifierUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryRulesVerifierUpdated)
				if err := _FaceRegistry.contract.UnpackLog(event, "RulesVerifierUpdated", log); err != nil {
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

// ParseRulesVerifierUpdated is a log parse operation binding the contract event 0x5a8b13378dd1cabe2f3617497f917ede72540af55149036f6afc92c031efbf1f.
//
// Solidity: event RulesVerifierUpdated(address oldVerifier, address newVerifier)
func (_FaceRegistry *FaceRegistryFilterer) ParseRulesVerifierUpdated(log types.Log) (*FaceRegistryRulesVerifierUpdated, error) {
	event := new(FaceRegistryRulesVerifierUpdated)
	if err := _FaceRegistry.contract.UnpackLog(event, "RulesVerifierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FaceRegistry contract.
type FaceRegistryUpgradedIterator struct {
	Event *FaceRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *FaceRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryUpgraded)
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
		it.Event = new(FaceRegistryUpgraded)
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
func (it *FaceRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryUpgraded represents a Upgraded event raised by the FaceRegistry contract.
type FaceRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FaceRegistry *FaceRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FaceRegistryUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FaceRegistryUpgradedIterator{contract: _FaceRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FaceRegistry *FaceRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FaceRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryUpgraded)
				if err := _FaceRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_FaceRegistry *FaceRegistryFilterer) ParseUpgraded(log types.Log) (*FaceRegistryUpgraded, error) {
	event := new(FaceRegistryUpgraded)
	if err := _FaceRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryUserRegisteredIterator is returned from FilterUserRegistered and is used to iterate over the raw logs and unpacked data for UserRegistered events raised by the FaceRegistry contract.
type FaceRegistryUserRegisteredIterator struct {
	Event *FaceRegistryUserRegistered // Event containing the contract specifics and raw log

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
func (it *FaceRegistryUserRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryUserRegistered)
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
		it.Event = new(FaceRegistryUserRegistered)
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
func (it *FaceRegistryUserRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryUserRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryUserRegistered represents a UserRegistered event raised by the FaceRegistry contract.
type FaceRegistryUserRegistered struct {
	UserAddress *big.Int
	FeatureHash *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUserRegistered is a free log retrieval operation binding the contract event 0x7c91c4b3b0da33c7349ab14f4fd11f582793495e5140eff2a1169f195cb81b8b.
//
// Solidity: event UserRegistered(uint256 userAddress, uint256 featureHash)
func (_FaceRegistry *FaceRegistryFilterer) FilterUserRegistered(opts *bind.FilterOpts) (*FaceRegistryUserRegisteredIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "UserRegistered")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryUserRegisteredIterator{contract: _FaceRegistry.contract, event: "UserRegistered", logs: logs, sub: sub}, nil
}

// WatchUserRegistered is a free log subscription operation binding the contract event 0x7c91c4b3b0da33c7349ab14f4fd11f582793495e5140eff2a1169f195cb81b8b.
//
// Solidity: event UserRegistered(uint256 userAddress, uint256 featureHash)
func (_FaceRegistry *FaceRegistryFilterer) WatchUserRegistered(opts *bind.WatchOpts, sink chan<- *FaceRegistryUserRegistered) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "UserRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryUserRegistered)
				if err := _FaceRegistry.contract.UnpackLog(event, "UserRegistered", log); err != nil {
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

// ParseUserRegistered is a log parse operation binding the contract event 0x7c91c4b3b0da33c7349ab14f4fd11f582793495e5140eff2a1169f195cb81b8b.
//
// Solidity: event UserRegistered(uint256 userAddress, uint256 featureHash)
func (_FaceRegistry *FaceRegistryFilterer) ParseUserRegistered(log types.Log) (*FaceRegistryUserRegistered, error) {
	event := new(FaceRegistryUserRegistered)
	if err := _FaceRegistry.contract.UnpackLog(event, "UserRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FaceRegistryVerifierUpdatedIterator is returned from FilterVerifierUpdated and is used to iterate over the raw logs and unpacked data for VerifierUpdated events raised by the FaceRegistry contract.
type FaceRegistryVerifierUpdatedIterator struct {
	Event *FaceRegistryVerifierUpdated // Event containing the contract specifics and raw log

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
func (it *FaceRegistryVerifierUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FaceRegistryVerifierUpdated)
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
		it.Event = new(FaceRegistryVerifierUpdated)
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
func (it *FaceRegistryVerifierUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FaceRegistryVerifierUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FaceRegistryVerifierUpdated represents a VerifierUpdated event raised by the FaceRegistry contract.
type FaceRegistryVerifierUpdated struct {
	OldVerifier common.Address
	NewVerifier common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVerifierUpdated is a free log retrieval operation binding the contract event 0x0243549a92b2412f7a3caf7a2e56d65b8821b91345363faa5f57195384065fcc.
//
// Solidity: event VerifierUpdated(address oldVerifier, address newVerifier)
func (_FaceRegistry *FaceRegistryFilterer) FilterVerifierUpdated(opts *bind.FilterOpts) (*FaceRegistryVerifierUpdatedIterator, error) {

	logs, sub, err := _FaceRegistry.contract.FilterLogs(opts, "VerifierUpdated")
	if err != nil {
		return nil, err
	}
	return &FaceRegistryVerifierUpdatedIterator{contract: _FaceRegistry.contract, event: "VerifierUpdated", logs: logs, sub: sub}, nil
}

// WatchVerifierUpdated is a free log subscription operation binding the contract event 0x0243549a92b2412f7a3caf7a2e56d65b8821b91345363faa5f57195384065fcc.
//
// Solidity: event VerifierUpdated(address oldVerifier, address newVerifier)
func (_FaceRegistry *FaceRegistryFilterer) WatchVerifierUpdated(opts *bind.WatchOpts, sink chan<- *FaceRegistryVerifierUpdated) (event.Subscription, error) {

	logs, sub, err := _FaceRegistry.contract.WatchLogs(opts, "VerifierUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FaceRegistryVerifierUpdated)
				if err := _FaceRegistry.contract.UnpackLog(event, "VerifierUpdated", log); err != nil {
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

// ParseVerifierUpdated is a log parse operation binding the contract event 0x0243549a92b2412f7a3caf7a2e56d65b8821b91345363faa5f57195384065fcc.
//
// Solidity: event VerifierUpdated(address oldVerifier, address newVerifier)
func (_FaceRegistry *FaceRegistryFilterer) ParseVerifierUpdated(log types.Log) (*FaceRegistryVerifierUpdated, error) {
	event := new(FaceRegistryVerifierUpdated)
	if err := _FaceRegistry.contract.UnpackLog(event, "VerifierUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
