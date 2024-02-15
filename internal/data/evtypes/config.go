package evtypes

import (
	"fmt"
	"time"

	"github.com/rarimo/rarime-points-svc/resources"
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
			Types []struct {
				Name        string     `fig:"name,required"`
				Description string     `fig:"description,required"`
				Reward      int64      `fig:"reward,required"`
				Title       string     `fig:"title,required"`
				Frequency   Frequency  `fig:"frequency,required"`
				ExpiresAt   *time.Time `fig:"expires_at"`
				NoAutoOpen  bool       `fig:"no_auto_open"`
				Disabled    bool       `fig:"disabled"`
			} `fig:"types,required"`
		}

		err := figure.Out(&raw).
			From(kv.MustGetStringMap(c.getter, "event_types")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out event_types: %s", err))
		}

		m := make(map[string]resources.EventStaticMeta, len(raw.Types))
		list := make([]resources.EventStaticMeta, 0, len(raw.Types))
		for _, t := range raw.Types {
			if !checkFreqValue(t.Frequency) {
				panic(fmt.Errorf("invalid frequency: %s", t.Frequency))
			}
			if t.Disabled {
				continue
			}

			meta := resources.EventStaticMeta{
				Name:        t.Name,
				Description: t.Description,
				Reward:      t.Reward,
				Title:       t.Title,
				Frequency:   t.Frequency.String(),
				ExpiresAt:   t.ExpiresAt,
				NoAutoOpen:  t.NoAutoOpen,
			}

			m[t.Name] = meta
			list = append(list, meta)
		}

		return Types{m, list}
	}).(Types)
}

func checkFreqValue(f Frequency) bool {
	switch f {
	case OneTime, Daily, Weekly, Unlimited:
		return true
	}
	return false
}
