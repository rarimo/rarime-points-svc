package sbtcheck

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rarimo/rarime-points-svc/internal/sbtcheck/erc721"
)

var ErrUnsupportedNetwork = errors.New("unsupported network")

type Connector struct {
	networks map[string]network
}

type network struct {
	caller  *erc721.IERC721Caller
	timeout time.Duration
}

func (c *Connector) IsSbtOwner(ctx context.Context, network, userAddress string) (bool, error) {
	net, ok := c.networks[network]
	if !ok {
		return false, fmt.Errorf("%w: %s", ErrUnsupportedNetwork, network)
	}

	toCtx, cancel := context.WithTimeout(ctx, net.timeout)
	defer cancel()

	balance, err := net.caller.BalanceOf(&bind.CallOpts{Context: toCtx}, common.HexToAddress(userAddress))
	if err != nil {
		return false, fmt.Errorf("check ERC721 SBT balance (network=%s userAddress=%s]: %w", network, userAddress, err)
	}

	return balance.Cmp(big.NewInt(0)) == 1, nil
}

func (c *Connector) IsNetworkSupported(network string) bool {
	_, ok := c.networks[network]
	return ok
}
