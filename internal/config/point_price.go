package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type PointsPrice struct {
	PointPriceURMO int64 `fig:"point_price_urmo,required"`
	Disabled       bool  `fig:"disabled"`
}

func (c *config) PointPrice() PointsPrice {
	return c.pointPrice.Do(func() interface{} {
		var cfg PointsPrice

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "withdrawal")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out withdrawal point price: %w", err))
		}

		return cfg
	}).(PointsPrice)
}
