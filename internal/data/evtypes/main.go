package evtypes

import (
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
)

const (
	TypeGetPoH           = "get_poh"
	TypeFreeWeekly       = "free_weekly"
	TypeBeReferred       = "be_referred"
	TypeReferralSpecific = "referral_specific"
)

type Types struct {
	m    map[string]resources.EventStaticMeta
	list []resources.EventStaticMeta
}

func (t Types) Get(name string, filters ...filter) *resources.EventStaticMeta {
	t.ensureInitialized()
	v, ok := t.m[name]
	if !ok || isFiltered(v, filters...) {
		return nil
	}

	return &v
}

func (t Types) List(filters ...filter) []resources.EventStaticMeta {
	t.ensureInitialized()
	res := make([]resources.EventStaticMeta, 0, len(t.list))
	for _, v := range t.list {
		if isFiltered(v, filters...) {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (t Types) Names(filters ...filter) []string {
	t.ensureInitialized()
	res := make([]string, 0, len(t.list))
	for _, v := range t.list {
		if isFiltered(v, filters...) {
			continue
		}
		res = append(res, v.Name)
	}
	return res
}

func (t Types) PrepareEvents(userDID string, filters ...filter) []data.Event {
	t.ensureInitialized()
	const extraCap = 1 // in case we append to the resulting slice outside the function
	events := make([]data.Event, 0, len(t.list)+extraCap)

	for _, et := range t.list {
		if isFiltered(et, filters...) {
			continue
		}

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

func (t Types) ensureInitialized() {
	if t.m == nil || t.list == nil {
		panic("event types are not correctly initialized")
	}
}
