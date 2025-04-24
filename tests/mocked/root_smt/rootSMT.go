package rootsmt

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RootSMTFilterer interface {
	FilterRootUpdated(opts *bind.FilterOpts, roots [][32]byte) (RootUpdatedIterator, error)
}

type RootUpdatedIterator interface {
	Next() bool
}

type RootSMTFiltererMock struct{}

func NewRootSMTFiltererMock(addr common.Address, client *ethclient.Client) (RootSMTFilterer, error) {
	return &RootSMTFiltererMock{}, nil
}

func (f *RootSMTFiltererMock) FilterRootUpdated(opts *bind.FilterOpts, roots [][32]byte) (RootUpdatedIterator, error) {
	return &MockRootUpdatedIterator{}, nil
}

type MockRootUpdatedIterator struct{}

func (i *MockRootUpdatedIterator) Next() bool {
	return true
}
