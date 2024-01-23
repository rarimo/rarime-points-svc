package evtypes

import (
	"errors"

	"github.com/rarimo/points-svc/resources"
)

var ErrInvalidID = errors.New("wrong event type id")

type Types struct {
	inner []resources.EventStaticMeta
}

func (t *Types) Get(id int) (resources.EventStaticMeta, error) {
	if t.inner == nil {
		panic("event types are not correctly initialized")
	}

	if id <= 0 || id > len(t.inner) {
		return resources.EventStaticMeta{}, ErrInvalidID
	}

	return t.inner[id-1], nil
}
