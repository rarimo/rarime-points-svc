package evtypes

import (
	"time"

	"github.com/rarimo/points-svc/resources"
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
