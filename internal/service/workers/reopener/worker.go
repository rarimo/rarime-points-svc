package reopener

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

type worker struct {
	name  string
	freq  evtypes.Frequency
	q     data.EventsQ
	types evtypes.Types
	log   *logan.Entry
}

func newWorker(cfg config.Config, freq evtypes.Frequency) *worker {
	name := fmt.Sprintf("reopener[%s]", freq.String())
	return &worker{
		name:  name,
		freq:  freq,
		q:     pg.NewEvents(cfg.DB().Clone()),
		types: cfg.EventTypes(),
		log:   cfg.Log().WithField("who", name),
	}
}

func (w *worker) job(ctx context.Context) {
	// types might expire, so it's required to get them before each run
	types := w.types.Names(evtypes.FilterByFrequency(w.freq), evtypes.FilterInactive)
	if len(types) == 0 {
		w.log.Info("No events to reopen: all types expired or no types with frequency exist")
		return
	}
	w.log.WithField("event_types", types).Debug("Reopening claimed events")

	running.WithThreshold(ctx, w.log, w.name, func(context.Context) (bool, error) {
		if err := w.reopenEvents(types); err != nil {
			return false, fmt.Errorf("reopen events: %w", err)
		}
		return true, nil
	}, retryPeriod, retryPeriod, 12)
}

func (w *worker) reopenEvents(types []string) error {
	log := w.log.WithField("event_types", types)

	events, err := w.q.New().FilterByType(types...).SelectReopenable()
	if err != nil {
		return fmt.Errorf("select reopenable events [types=%v]: %w", types, err)
	}
	if len(events) == 0 {
		log.Info("No events to reopen: no claimed events found for provided types")
		return nil
	}
	log.Infof("%d (DID, type) pairs to reopen: %v", len(events), events)

	err = w.q.New().Insert(prepareForReopening(events)...)
	if err != nil {
		return fmt.Errorf("insert events for reopening: %w", err)
	}

	w.log.Infof("Reopened %d events", len(events))
	return nil
}
