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
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

const retryPeriod = 5 * time.Minute

type worker struct {
	name  string
	freq  evtypes.Frequency
	types evtypes.Types
	q     data.EventsQ
	log   *logan.Entry
}

func Run(ctx context.Context, cfg config.Config) {
	var (
		atUTC  = gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))
		daily  = newWorker(cfg, evtypes.Daily)
		weekly = newWorker(cfg, evtypes.Weekly)
	)
	if err := daily.initialRun(); err != nil {
		panic(fmt.Errorf("failed to do initial run for daily events: %w", err))
	}
	if err := weekly.initialRun(); err != nil {
		panic(fmt.Errorf("failed to do initial run for weekly events: %w", err))
	}

	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.UTC),
		gocron.WithLogger(newLogger(cfg.Log())),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize scheduler: %w", err))
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, atUTC),
		gocron.NewTask(daily.job, ctx),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize daily job: %w", err))
	}
	_, err = scheduler.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(time.Monday), atUTC),
		gocron.NewTask(weekly.job, ctx),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize weekly job: %w", err))
	}

	scheduler.Start()
	<-ctx.Done()
	if err = scheduler.Shutdown(); err != nil {
		cfg.Log().WithError(err).Error("Scheduler shutdown failed")
		return
	}
	cfg.Log().Info("Scheduler shutdown succeeded")
}

func newWorker(cfg config.Config, freq evtypes.Frequency) *worker {
	name := fmt.Sprintf("reopener[%s]", freq.String())
	return &worker{
		name:  name,
		freq:  freq,
		types: cfg.EventTypes(),
		q:     pg.NewEvents(cfg.DB().Clone()),
		log:   cfg.Log().WithField("who", name),
	}
}

func (w *worker) job(ctx context.Context) {
	// types might expire, so it's required to get them before each run
	types := w.types.NamesByFrequency(w.freq)
	if len(types) == 0 {
		w.log.Info("No events to reopen: all types expired or no types with frequency exist")
		return
	}
	w.log.WithField("event_types", types).
		Debug("Reopening claimed events")

	running.WithThreshold(ctx, w.log, w.name, func(context.Context) (bool, error) {
		if err := w.reopenEvents(types, false); err != nil {
			return false, fmt.Errorf("reopen events: %w", err)
		}
		return true, nil
	}, retryPeriod, retryPeriod, 12)
}

func (w *worker) reopenEvents(types []string, initRun bool) error {
	q := w.q.New().FilterByType(types...)

	if initRun {
		filter := w.beforeTimeFilter()
		w.log.WithField("event_types", types).
			Debugf("Reopening claimed events before %d", filter)
		q.FilterByUpdatedAtBefore(filter)
	}

	events, err := q.SelectReopenable()
	if err != nil {
		return fmt.Errorf("select reopenable events [types=%v]: %w", types, err)
	}
	if len(events) == 0 {
		w.log.Info("No events to reopen")
		return nil
	}

	err = w.q.New().Insert(prepareForReopening(events)...)
	if err != nil {
		return fmt.Errorf("insert events for reopening: %w", err)
	}

	w.log.Infof("Reopened %d events", len(events))
	return nil
}

func prepareForReopening(events []data.ReopenableEvent) []data.Event {
	res := make([]data.Event, len(events))

	for i, ev := range events {
		res[i] = data.Event{
			UserDID: ev.UserDID,
			Type:    ev.Type,
			Status:  data.EventOpen,
		}

		if ev.Type == evtypes.TypeFreeWeekly {
			res[i].Status = data.EventFulfilled
		}
	}

	return res
}
