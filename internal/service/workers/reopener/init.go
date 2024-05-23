package reopener

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/workers/cron"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

func initialRun(cfg config.Config) error {
	var (
		q   = pg.NewEvents(cfg.DB().Clone())
		log = cfg.Log().WithField("who", "reopener[initializer]")
		col = &initCollector{
			q:     q,
			types: cfg.EventTypes(),
			log:   log,
		}
	)

	events, err := col.collect()
	if err != nil {
		return fmt.Errorf("collect events: %w", err)
	}

	err = q.New().Insert(prepareForReopening(events)...)
	if err != nil {
		return fmt.Errorf("insert events to be opened: %w", err)
	}

	log.Infof("Reopened %d events on the initial run", len(events))
	return nil
}

type initCollector struct {
	q     data.EventsQ
	types evtypes.Types
	log   *logan.Entry
}

func (c *initCollector) collect() ([]data.ReopenableEvent, error) {
	var (
		now       = time.Now().UTC()
		monOffset = int(time.Monday - now.Weekday())
		midnight  = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		weekStart = midnight.AddDate(0, 0, monOffset).Unix()
	)

	daily, err := c.selectReopenable(evtypes.Daily, midnight.Unix())
	if err != nil {
		return nil, fmt.Errorf("select daily events: %w", err)
	}

	weekly, err := c.selectReopenable(evtypes.Weekly, weekStart)
	if err != nil {
		return nil, fmt.Errorf("select weekly events: %w", err)
	}

	absent, err := c.selectAbsent()
	if err != nil {
		return nil, fmt.Errorf("select absent events: %w", err)
	}

	dw := append(daily, weekly...)
	return append(dw, absent...), nil
}

func (c *initCollector) selectReopenable(freq evtypes.Frequency, before int64) ([]data.ReopenableEvent, error) {
	types := c.types.Names(evtypes.FilterByFrequency(freq), evtypes.FilterInactive)

	if len(types) == 0 {
		return nil, nil
	}

	res, err := c.q.New().FilterByType(types...).
		FilterByUpdatedAtBefore(before).
		SelectReopenable()
	if err != nil {
		return nil, fmt.Errorf("select reopenable events [freq=%s before=%d types=%v]: %w", freq, before, types, err)
	}

	log := c.log.WithFields(logan.F{
		"frequency": freq,
		"before":    before,
		"types":     types,
	})

	if len(res) == 0 {
		log.Debug("No events to reopen on initial run")
		return nil, nil
	}

	log.Infof("%d (nullifier, type) pairs to reopen: %v", len(res), res)
	return res, nil
}

func (c *initCollector) selectAbsent() ([]data.ReopenableEvent, error) {
	types := c.types.Names(evtypes.FilterNotOpenable)
	if len(types) == 0 {
		c.log.Debug("No openable event types are active, skip absent types selection")
		return nil, nil
	}

	res, err := c.q.New().SelectAbsentTypes(types...)
	if err != nil {
		return nil, fmt.Errorf("select events with absent types [types=%v]: %w", types, err)
	}

	log := c.log.WithField("types", types)
	if len(res) == 0 {
		log.Debug("No new event types found to open for new users")
		return nil, nil
	}

	log.Infof("%d new (nullifier, type) pairs to open: %v", len(res), res)
	return res, nil
}

func runStartingWatchers(ctx context.Context, cfg config.Config) error {
	log := cfg.Log().WithField("who", "opener[initializer]")

	notStartedEv := cfg.EventTypes().List(func(ev evtypes.EventConfig) bool {
		return ev.Disabled || !evtypes.FilterNotStarted(ev) || evtypes.FilterExpired(ev)
	})

	if len(notStartedEv) == 0 {
		log.Info("No events to open at Start time: all types already opened or there no types with StartAt")
		return nil
	}

	for _, ev := range notStartedEv {
		_, err := cron.NewJob(
			gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(*ev.StartsAt)),
			gocron.NewTask(startingWatcher(cfg, ev.Name), ctx),
			gocron.WithName(fmt.Sprintf("opener[%s]", ev.Name)),
		)

		if err != nil {
			return fmt.Errorf("opener: failed to initialize job: %w", err)
		}
	}

	return nil
}

func startingWatcher(cfg config.Config, name string) func(context.Context) {
	return func(ctx context.Context) {
		log := cfg.Log().WithField("who", fmt.Sprintf("opener[%s]", name))

		var balances []data.Balance
		var err error

		running.WithThreshold(ctx, log, fmt.Sprintf("opener[%s]", name), func(context.Context) (bool, error) {
			if balances, err = pg.NewBalances(cfg.DB().Clone()).Select(); err != nil {
				return false, err
			}
			return true, nil
		}, retryPeriod, retryPeriod, maxRetries)

		events := make([]data.Event, len(balances))
		status := data.EventOpen
		if name == evtypes.TypeFreeWeekly {
			status = data.EventFulfilled
		}

		for i, balance := range balances {
			events[i] = data.Event{Nullifier: balance.Nullifier, Type: name, Status: status}
		}

		running.WithThreshold(ctx, log, fmt.Sprintf("opener[%s]", name), func(context.Context) (bool, error) {
			if err = pg.NewEvents(cfg.DB().Clone()).Insert(events...); err != nil {
				return false, err
			}
			return true, nil
		}, retryPeriod, retryPeriod, maxRetries)
	}
}
