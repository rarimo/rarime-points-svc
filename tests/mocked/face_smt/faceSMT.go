package facesmt

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type FaceSMTFilterer interface {
	FilterRootUpdated(opts *bind.FilterOpts, roots [][32]byte) (RootUpdatedIterator, error)
}

type RootUpdatedIterator interface {
	Next() bool
}

type FaceSMTFiltererMock struct{}

func NewFaceSMTFiltererMock(addr common.Address, client *ethclient.Client) (FaceSMTFilterer, error) {
	return &FaceSMTFiltererMock{}, nil
}

func (f *FaceSMTFiltererMock) FilterRootUpdated(opts *bind.FilterOpts, roots [][32]byte) (RootUpdatedIterator, error) {
	return &MockRootUpdatedIterator{}, nil
}

type MockRootUpdatedIterator struct{}

func (i *MockRootUpdatedIterator) Next() bool {
	return true
}
