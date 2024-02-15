package evtypes

import (
	"time"

	"github.com/rarimo/rarime-points-svc/resources"
)

type filter func(resources.EventStaticMeta) bool

func FilterExpired(ev resources.EventStaticMeta) bool {
	return ev.ExpiresAt != nil && ev.ExpiresAt.Before(time.Now().UTC())
}

func FilterNoAutoOpen(ev resources.EventStaticMeta) bool {
	return ev.NoAutoOpen
}

func FilterByFrequency(f Frequency) func(resources.EventStaticMeta) bool {
	return func(ev resources.EventStaticMeta) bool {
		return ev.Frequency != f.String()
	}
}

func isFiltered(ev resources.EventStaticMeta, filters ...filter) bool {
	for _, f := range filters {
		if f(ev) {
			return true
		}
	}
	return false
}
