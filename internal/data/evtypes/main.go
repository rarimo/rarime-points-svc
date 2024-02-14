package evtypes

import (
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

const (
	TypeGetPoH     = "get_poh"
	TypeFreeWeekly = "free_weekly"
	TypeBeReferred = "be_referred"
)

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

func (t Types) PrepareOpenEvents(userDID string) []data.Event {
	const extraCap = 1 // in case we append to the resulting slice outside the function
	evTypes := t.List()
	events := make([]data.Event, 0, len(evTypes)+extraCap)

	for _, et := range evTypes {
		status := data.EventOpen
		if et.Name == TypeFreeWeekly {
			status = data.EventFulfilled
		}

		events = append(events, data.Event{
			UserDID: userDID,
			Type:    et.Name,
			Status:  status,
		})
	}

	return events
}

// List returns non-expired and auto-opening event types
func (t Types) List() []resources.EventStaticMeta {
	if t.inner == nil {
		panic("event types are not correctly initialized")
	}

	res := make([]resources.EventStaticMeta, 0, len(t.inner))
	for _, v := range t.inner {
		if v.NoAutoOpen || isExpiredEvent(v) {
			continue
		}
		res = append(res, v)
	}

	return res
}

func (t Types) NamesByFrequency(f Frequency) []string {
	if t.inner == nil {
		panic("event types are not correctly initialized")
	}

	res := make([]string, 0, len(t.inner))
	for _, v := range t.inner {
		if v.Frequency != f.String() || isExpiredEvent(v) {
			continue
		}
		res = append(res, v.Name)
	}

	return res
}

func (t Types) IsExpired(name string) bool {
	evType := t.Get(name)
	if evType == nil {
		return false
	}

	return isExpiredEvent(*evType)
}

func isExpiredEvent(ev resources.EventStaticMeta) bool {
	return ev.ExpiresAt != nil && ev.ExpiresAt.Before(time.Now().UTC())
}
