package config

import (
	"fmt"
	"slices"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type Level struct {
	Level             int  `fig:"lvl,required"`
	Threshold         int  `fig:"threshold,required"`
	Referrals         int  `fig:"referrals,required"`
	WithdrawalAllowed bool `fig:"withdrawal_allowed"`
}

type Levels map[int]Level

func (c *config) Levels() Levels {
	return c.levels.Do(func() interface{} {
		var cfg struct {
			Lvls []Level `fig:"levels,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "levels")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out levels config: %w", err))
		}

		if len(cfg.Lvls) == 0 {
			panic(fmt.Errorf("no levels provided in config: %w", err))
		}

		res := make(Levels, len(cfg.Lvls))
		for _, v := range cfg.Lvls {
			res[v.Level] = v
		}

		return res
	}).(Levels)
}

// Calculate new lvl. New lvl always greater then current level
func (l Levels) LvlUp(currentLevel int, totalAmount int64) (refCoundToAdd int, newLevel int) {
	lvls := make([]int, 0, len(l))
	for k, v := range l {
		if k <= currentLevel {
			continue
		}
		if int64(v.Threshold) > totalAmount {
			break
		}

		refCoundToAdd += v.Referrals
		lvls = append(lvls, k)
	}

	if len(lvls) == 0 {
		return 0, currentLevel
	}

	newLevel = slices.Max(lvls)
	return
}

// slices.Min will not panic because of previous logic
func (l Levels) MinLvl() int {
	lvls := make([]int, 0, len(l))
	for k := range l {
		lvls = append(lvls, k)
	}

	return slices.Min(lvls)
}
