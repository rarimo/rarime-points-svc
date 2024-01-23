package evtypes

import (
	"github.com/rarimo/points-svc/resources"
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
