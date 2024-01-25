package evtypes

import (
	"database/sql"
	"time"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/resources"
)

type Frequency string

func (f Frequency) String() string {
	return string(f)
}

const (
	OneTime   Frequency = "one-time"
	Daily     Frequency = "daily"
	Weekly    Frequency = "weekly"
	Unlimited Frequency = "unlimited"
	Custom    Frequency = "custom"
)

const TypeGetPoH = "get_poh"

type Types struct {
	inner map[string]resources.EventStaticMeta
}

func (t Types) Get(name string) *resources.EventStaticMeta {
	if t.inner == nil {
		panic("event types are not correctly initialized")
	}

	v, ok := t.inner[name]
	if !ok {
		return nil
	}

	return &v
}

func (t Types) PrepareOpenEvents(balanceID string) []data.Event {
	evTypes := t.List()
	events := make([]data.Event, len(evTypes))

	for i, evType := range evTypes {
		// TODO: add advanced logic for specific event types
		// for example, proof verification events should appear after the proof was issued
		events[i] = data.Event{
			BalanceID: balanceID,
			Type:      evType.Name,
			Status:    data.EventOpen,
			PointsAmount: sql.NullInt32{
				Int32: evType.Reward,
				Valid: true,
			},
		}
	}

	return events
}

// List returns non-expired event types
func (t Types) List() []resources.EventStaticMeta {
	if t.inner == nil {
		panic("event types are not correctly initialized")
	}

	res := make([]resources.EventStaticMeta, 0, len(t.inner))
	for _, v := range t.inner {
		if isExpiredEvent(v) {
			continue
		}
		res = append(res, v)
	}

	return res
}

func isExpiredEvent(ev resources.EventStaticMeta) bool {
	return ev.ExpiresAt != nil && ev.ExpiresAt.Before(time.Now().UTC())
}
