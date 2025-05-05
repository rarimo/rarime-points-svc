package evtypes

import (
	"time"
)

// Filter functions work in the following way:
//
// 1. For FilterBy* functions, the config is only added when it matches the filter:
// FilterByName(name1, name2) will only return events with name1 or name2
//
// 2. For other Filter* functions, the configs matching the filter are excluded:
// FilterExpired eliminates all expired events (instead of including only them)

type filter func(EventConfig) bool

func FilterExpired(ev EventConfig) bool {
	return ev.ExpiresAt != nil && ev.ExpiresAt.Before(time.Now().UTC())
}

func FilterNotStarted(ev EventConfig) bool {
	return ev.StartsAt != nil && ev.StartsAt.After(time.Now().UTC())
}

func FilterInactive(ev EventConfig) bool {
	return ev.Disabled || FilterExpired(ev) || FilterNotStarted(ev)
}

func FilterNotOpenable(ev EventConfig) bool {
	return FilterInactive(ev) || ev.NoAutoOpen
}

func FilterByFrequency(f Frequency) func(EventConfig) bool {
	return func(ev EventConfig) bool {
		return ev.Frequency != f
	}
}

func FilterByNames(names ...string) func(EventConfig) bool {
	return func(ev EventConfig) bool {
		if len(names) == 0 {
			return false
		}
		for _, name := range names {
			if ev.Name == name {
				return false
			}
		}
		return true
	}
}

func FilterByFlags(flags ...string) func(EventConfig) bool {
	return func(ev EventConfig) bool {
		if len(flags) == 0 {
			return false
		}
		for _, flag := range flags {
			if ev.Flag() == flag {
				return false
			}
		}
		return true
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
