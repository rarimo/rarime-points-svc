package config

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type Maintenance struct {
	IsMaintenance bool `fig:"is_maintenance"`
}

func (c *config) Maintenance() Maintenance {
	return c.maintenance.Do(func() interface{} {
		var cfg Maintenance

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "maintenance")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out is_maintenance: %w", err))
		}

		return cfg
	}).(Maintenance)
}
