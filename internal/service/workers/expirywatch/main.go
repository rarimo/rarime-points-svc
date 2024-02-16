package expirywatch

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/cron"
)

const retryPeriod = 1 * time.Minute

func Run(ctx context.Context, cfg config.Config) {
	w := newWatcher(cfg)
	if err := w.initialRun(); err != nil {
		panic(fmt.Errorf("expiry-watcher: initial run failed: %w", err))
	}

	cron.Init(cfg.Log())
	expirable := w.types.List(func(ev evtypes.EventConfig) bool {
		return ev.Disabled || ev.ExpiresAt == nil || evtypes.FilterExpired(ev)
	})

	for _, ev := range expirable {
		if ev.ExpiresAt.Before(time.Now().UTC()) {
			continue // although we filtered expired, ensure extra safety due to possible delay
		}

		_, err := cron.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(*ev.ExpiresAt)),
			gocron.NewTask(w.job, ctx, ev.Name),
			gocron.WithName(fmt.Sprintf("expiry-watch[%s]", ev.Name)),
		)
		if err != nil {
			panic(fmt.Errorf("failed to initialize job [event_type=%+v]: %w", ev, err))
		}
	}

	cron.Start(ctx)
}
