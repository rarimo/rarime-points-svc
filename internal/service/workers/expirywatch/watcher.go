package expirywatch

import (
	"context"
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

type watcher struct {
	q     data.EventsQ
	types evtypes.Types
	log   *logan.Entry
}

func newWatcher(cfg config.Config) *watcher {
	return &watcher{
		q:     pg.NewEvents(cfg.DB().Clone()),
		types: cfg.EventTypes(),
		log:   cfg.Log().WithField("who", "expiry-watch"),
	}
}

func (w *watcher) initialRun() error {
	expired := w.types.Names(func(ev evtypes.EventConfig) bool {
		return !ev.Disabled && !evtypes.FilterExpired(ev)
	})

	if len(expired) == 0 {
		w.log.Debug("No events were disabled or have expired")
		return nil
	}

	return w.cleanOpen(expired...)
}

func (w *watcher) cleanOpen(types ...string) error {
	deleted, err := w.q.New().FilterByType(types...).FilterByStatus(data.EventOpen).Delete()
	if err != nil {
		return fmt.Errorf("clean open events [types=%v]: %w", types, err)
	}

	w.log.Infof("Deleted %d expired and disabled open events, types: %v", deleted, types)
	return nil
}

func (w *watcher) job(ctx context.Context, evType string) {
	name := fmt.Sprintf("expiry-watch[%s]", evType)
	log := w.log.WithField("who", name)

	running.WithThreshold(ctx, log, name, func(context.Context) (bool, error) {
		if err := w.cleanOpen(evType); err != nil {
			return false, fmt.Errorf("clean open events: %w", err)
		}
		return true, nil
	}, retryPeriod, retryPeriod, 12)
}
