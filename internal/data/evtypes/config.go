package evtypes

import (
	"fmt"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type EventTypeser interface {
	EventTypes() Types
}

type config struct {
	once   comfig.Once
	getter kv.Getter
}

func NewConfig(getter kv.Getter) EventTypeser {
	return &config{getter: getter}
}

func (c *config) EventTypes() Types {
	return c.once.Do(func() interface{} {
		var raw struct {
			Types           []EventConfig  `fig:"types,required"`
			PassportRewards map[string]int `fig:"passport_rewards"`
		}

		err := figure.Out(&raw).
			From(kv.MustGetStringMap(c.getter, "event_types")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out event_types: %s", err))
		}

		m := make(map[string]EventConfig, len(raw.Types))
		for _, t := range raw.Types {
			if !checkFreqValue(t.Frequency) {
				panic(fmt.Errorf("invalid frequency: %s", t.Frequency))
			}
			m[t.Name] = t
		}

		if _, ok := m[TypePassportScan]; ok {
			if _, ok := raw.PassportRewards[PassportRewardRequiredAge]; !ok {
				panic(fmt.Errorf("absent required field: %s", PassportRewardRequiredAge))
			} else if _, ok := raw.PassportRewards[PassportRewardRequiredNationality]; !ok {
				panic(fmt.Errorf("absent required field: %s", PassportRewardRequiredNationality))
			}
			return Types{m, raw.Types, raw.PassportRewards}
		}

		if len(raw.PassportRewards) != 0 {
			panic(fmt.Errorf("rewards exists, but event PassportScan not exists"))
		}
		return Types{m, raw.Types, raw.PassportRewards}
	}).(Types)
}

func checkFreqValue(f Frequency) bool {
	switch f {
	case OneTime, Daily, Weekly, Unlimited:
		return true
	}
	return false
}
