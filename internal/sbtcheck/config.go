package sbtcheck

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rarimo/rarime-points-svc/internal/sbtcheck/erc721"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const baseTimeout = 5 * time.Second

type SbtChecker interface {
	SbtCheck() *Connector
}

type config struct {
	once   comfig.Once
	getter kv.Getter
}

func NewConfig(getter kv.Getter) SbtChecker {
	return &config{getter: getter}
}

func (c *config) SbtCheck() *Connector {
	return c.once.Do(func() interface{} {
		var cfg struct {
			Networks []struct {
				Name           string        `fig:"name,required"`
				RPC            string        `fig:"rpc,required"`
				Contract       string        `fig:"contract,required"`
				RequestTimeout time.Duration `fig:"request_timeout"`
			} `fig:"networks,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "sbt_check")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out sbt_check: %s", err))
		}

		nmap := make(map[string]network, len(cfg.Networks))
		for _, net := range cfg.Networks {

			cli, err := ethclient.Dial(net.RPC)
			if err != nil {
				panic(fmt.Errorf("failed to connect to rpc: %w", err))
			}

			caller, err := erc721.NewIERC721Caller(common.HexToAddress(net.Contract), cli)
			if err != nil {
				panic(fmt.Errorf("failed to init contract caller: %w", err))
			}

			if net.RequestTimeout == 0 {
				net.RequestTimeout = baseTimeout
			}

			nmap[net.Name] = network{
				caller:  caller,
				timeout: net.RequestTimeout,
			}
		}

		return &Connector{networks: nmap}
	}).(*Connector)
}
