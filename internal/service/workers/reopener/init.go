package reopener

import (
	"fmt"
	"time"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
)

func (w *worker) initialRun() error {
	types := w.types.NamesByFrequency(w.freq)
	if len(types) == 0 {
		w.log.Info("No events to reopen: all types expired or no types with frequency exist")
		return nil
	}

	filter := w.beforeTimeFilter()
	w.log.WithField("event_types", types).
		Debugf("Reopening claimed events before %s", time.Unix(filter, 0).UTC())

	count, err := w.q.New().
		FilterByType(types...).
		FilterByStatus(data.EventClaimed).
		FilterByUpdatedAtBefore(filter).
		Reopen()

	if err != nil {
		return fmt.Errorf("reopen events: %w", err)
	}

	w.log.Infof("Reopened %d events on initial run", count)
	return nil
}

func (w *worker) beforeTimeFilter() int64 {
	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch w.freq {
	case evtypes.Daily:
		return midnight.Unix()
	case evtypes.Weekly:
		// current_day + (monday - current_day) = monday
		offset := int(time.Monday - now.Weekday())
		return midnight.AddDate(0, 0, offset).Unix()
	default:
		panic(fmt.Errorf("unexpected frequency: %s", w.freq))
	}
}
