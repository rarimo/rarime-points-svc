package evtypes

import (
	"time"
)

type filter func(EventConfig) bool

func FilterExpired(ev EventConfig) bool {
	return ev.ExpiresAt != nil && ev.ExpiresAt.Before(time.Now().UTC())
}

func FilterInactive(ev EventConfig) bool {
	return ev.Disabled || FilterExpired(ev)
}

func FilterNotOpenable(ev EventConfig) bool {
	return FilterInactive(ev) || ev.NoAutoOpen
}

func FilterByFrequency(f Frequency) func(EventConfig) bool {
	return func(ev EventConfig) bool {
		return ev.Frequency != f
	}
}

func isFiltered(ev EventConfig, filters ...filter) bool {
	for _, f := range filters {
		if f(ev) {
			return true
		}
	}
	return false
}
