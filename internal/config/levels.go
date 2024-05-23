package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Level struct {
	Level             int  `fig:"lvl,required"`
	Threshold         int  `fig:"threshold,required"`
	Referrals         int  `fig:"referrals,required"`
	WithdrawalAllowed bool `fig:"withdrawal_allowed"`
}

type Levels map[int]Level

func (c *config) Leveler() Levels {
	return c.leveler.Do(func() interface{} {
		var cfg struct {
			Lvls []Level `fig:"levels,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "levels")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out levels config"))
		}

		res := make(Levels)
		for _, v := range cfg.Lvls {
			res[v.Level] = v
		}

		return res
	}).(Levels)
}
