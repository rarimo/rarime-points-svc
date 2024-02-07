package reopener

import (
	"fmt"
	"time"

	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
)

func (w *worker) initialRun() error {
	types := w.types.NamesByFrequency(w.freq)
	if len(types) == 0 {
		w.log.Info("Initial run: no events to reopen: all types expired or no types with frequency exist")
		return nil
	}
	w.log.WithField("event_types", types).
		Debug("Initial run: reopening claimed/reserved events")

	return w.reopenEvents(types, true)
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
