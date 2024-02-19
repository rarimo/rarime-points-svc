package reopener

import (
	"fmt"
	"time"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"gitlab.com/distributed_lab/logan/v3"
)

func initialRun(cfg config.Config) error {
	var (
		q   = pg.NewEvents(cfg.DB())
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

	log.Infof("%d (DID, type) pairs to reopen: %v", len(res), res)
	return res, nil
}

func (c *initCollector) selectAbsent() ([]data.ReopenableEvent, error) {
	types := c.types.Names(evtypes.FilterNotOpenable)

	res, err := c.q.New().SelectAbsentTypes(types...)
	if err != nil {
		return nil, fmt.Errorf("select events with absent types [types=%v]: %w", types, err)
	}

	log := c.log.WithField("types", types)
	if len(res) == 0 {
		log.Debug("No new event types found to open for new users")
		return nil, nil
	}

	log.Infof("%d new (DID, type) pairs to open: %v", len(res), res)
	return res, nil
}
