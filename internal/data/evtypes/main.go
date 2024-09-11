package evtypes

import (
	"net/url"
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
)

const (
	TypeGetPoH           = "get_poh"
	TypeFreeWeekly       = "free_weekly"
	TypeBeReferred       = "be_referred"
	TypeReferralSpecific = "referral_specific"
	TypePassportScan     = "passport_scan"
	TypeEarlyTest        = "early_test"
)

const (
	FlagActive     = "active"
	FlagNotStarted = "not_started"
	FlagExpired    = "expired"
	FlagDisabled   = "disabled"
)

type EventConfig struct {
	Name             string     `fig:"name,required"`
	Description      string     `fig:"description,required"`
	ShortDescription string     `fig:"short_description,required"`
	Reward           int64      `fig:"reward,required"`
	Title            string     `fig:"title,required"`
	Frequency        Frequency  `fig:"frequency,required"`
	StartsAt         *time.Time `fig:"starts_at"`
	ExpiresAt        *time.Time `fig:"expires_at"`
	NoAutoOpen       bool       `fig:"no_auto_open"`
	AutoClaim        bool       `fig:"auto_claim"`
	Disabled         bool       `fig:"disabled"`
	ActionURL        *url.URL   `fig:"action_url"`
	Logo             *url.URL   `fig:"logo"`
}

func (e EventConfig) Flag() string {
	switch {
	case e.Disabled:
		return FlagDisabled
	case FilterNotStarted(e):
		return FlagNotStarted
	case FilterExpired(e):
		return FlagExpired
	default:
		return FlagActive
	}
}

func (e EventConfig) Resource() resources.EventStaticMeta {
	safeConv := func(u *url.URL) *string {
		if u == nil {
			return nil
		}
		s := u.String()
		return &s
	}

	return resources.EventStaticMeta{
		Name:             e.Name,
		Description:      e.Description,
		ShortDescription: e.ShortDescription,
		Reward:           e.Reward,
		Title:            e.Title,
		Frequency:        e.Frequency.String(),
		StartsAt:         e.StartsAt,
		ExpiresAt:        e.ExpiresAt,
		ActionUrl:        safeConv(e.ActionURL),
		Logo:             safeConv(e.Logo),
		Flag:             e.Flag(),
	}
}

type Types struct {
	m    map[string]EventConfig
	list []EventConfig
}

func (t Types) Get(name string, filters ...filter) *EventConfig {
	t.ensureInitialized()
	v, ok := t.m[name]
	if !ok || isFiltered(v, filters...) {
		return nil
	}

	return &v
}

func (t Types) List(filters ...filter) []EventConfig {
	t.ensureInitialized()
	res := make([]EventConfig, 0, len(t.list))
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

func (t Types) PrepareEvents(nullifier string, filters ...filter) []data.Event {
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
			Nullifier: nullifier,
			Type:      et.Name,
			Status:    status,
		})
	}

	return events
}

func (t Types) ensureInitialized() {
	if t.m == nil || t.list == nil {
		panic("event types are not correctly initialized")
	}
}
