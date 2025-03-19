package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type ExpiredCode struct {
	Code string `fig:"code"`
}

func (c *config) ExpiredCode() ExpiredCode {
	return c.expiredCode.Do(func() interface{} {
		var cfg ExpiredCode

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "expired_code")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out expired_code: %w", err))
		}

		return cfg
	}).(ExpiredCode)
}
