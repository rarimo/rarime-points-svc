package config

import (
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Countrier interface {
	Countries() Countries
}

type Countries map[string]Country

type Country struct {
	Code              string `fig:"code,required"`
	ReserveLimit      int64  `fig:"reserve_limit,required"`
	ReserveAllowed    bool   `fig:"reserve_allowed,required"`
	WithdrawalAllowed bool   `fig:"withdrawal_allowed,required"`
}

type countriesCfg struct {
	once   comfig.Once
	getter kv.Getter
}

func NewCountrier(getter kv.Getter) Countrier {
	return &countriesCfg{getter: getter}
}

func (c *countriesCfg) Countries() Countries {
	return c.once.Do(func() interface{} {
		var cfg struct {
			Countries []Country `fig:"countries,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "countries")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out countries: %s", err))
		}

		countries := make(Countries, len(cfg.Countries))
		for _, country := range cfg.Countries {
			if country.ReserveLimit < 0 {
				panic(fmt.Errorf("reserve limit can't be less than 0 [code=%s]", country.Code))
			}

			countries[country.Code] = country
		}

		if _, ok := countries[data.DefaultCountryCode]; !ok {
			panic("there no default country, default country must have code = 'default'")
		}

		return countries
	}).(Countries)
}
