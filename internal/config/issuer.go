package config

import (
	"fmt"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type IssuerConfig struct {
	Host             string `fig:"host,required"`
	Username         string `fig:"username,required"`
	Password         string `fig:"password,required"`
	CredentialSchema string `fig:"schema, required"`
	Type             string `fig:"type"`
}

func (c *config) IssuerConfig() *IssuerConfig {
	return c.issuer.Do(func() interface{} {
		var cfg IssuerConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "issuer")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out withdrawal point price: %w", err))
		}

		return &cfg
	}).(*IssuerConfig)
}
