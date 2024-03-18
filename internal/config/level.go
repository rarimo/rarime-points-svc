package config

import (
	"fmt"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

func (c *config) Levels() []int64 {
	return c.levels.Do(func() interface{} {
		var cfg struct {
			Levels []int64 `fig:"levels,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "levels")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out withdrawal point price: %w", err))
		}

		return cfg.Levels
	}).([]int64)
}
