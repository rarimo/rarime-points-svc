package evtypes

import (
	"fmt"
	"time"

	"github.com/rarimo/points-svc/resources"
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
				Reward      int32      `fig:"reward,required"`
				Title       string     `fig:"title,required"`
				ExpiresAt   *time.Time `fig:"expires_at"`
			} `fig:"types,required"`
		}

		err := figure.Out(&raw).
			From(kv.MustGetStringMap(c.getter, "event_types")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out event_types: %s", err))
		}

		inner := make(map[string]resources.EventStaticMeta, len(raw.Types))
		for _, t := range raw.Types {
			inner[t.Name] = resources.EventStaticMeta{
				Name:        t.Name,
				Description: t.Description,
				Reward:      t.Reward,
				Title:       t.Title,
				ExpiresAt:   t.ExpiresAt,
			}
		}

		return Types{inner}
	}).(Types)
}
