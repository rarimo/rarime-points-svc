package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

func (c *config) PointPrice() uint64 {
	return c.pointPrice.Do(func() interface{} {
		var cfg struct {
			PointPriceURMO uint64 `fig:"point_price_urmo,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "withdrawal")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out withdrawal point price: %w", err))
		}

		return cfg.PointPriceURMO
	}).(uint64)
}
