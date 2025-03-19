package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type ExterminatedCode struct {
	Code string `fig:"code"`
}

func (c *config) ExterminatedCode() ExterminatedCode {
	return c.exterminatedCode.Do(func() interface{} {
		var cfg ExterminatedCode

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "exterminated_code")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out exterminated_code: %w", err))
		}

		return cfg
	}).(ExterminatedCode)
}
