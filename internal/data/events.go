package data

import (
	"database/sql"
	"encoding/json"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type EventStatus string

const (
	EventOpen      EventStatus = "open"
	EventFulfilled EventStatus = "fulfilled"
	EventClaimed   EventStatus = "claimed"
)

func (s EventStatus) String() string {
	return string(s)
}

type Event struct {
	ID           string         `db:"id"`
	Nullifier    string         `db:"nullifier"`
	Type         string         `db:"type"`
	Status       EventStatus    `db:"status"`
	CreatedAt    int32          `db:"created_at"`
	UpdatedAt    int32          `db:"updated_at"`
	Meta         Jsonb          `db:"meta"`
	PointsAmount *int64         `db:"points_amount"`
	ExternalID   sql.NullString `db:"external_id"` // hidden from client
}

// ReopenableEvent is a pair that is sufficient to build a new open event with a specific type for a user
type ReopenableEvent struct {
	Nullifier string `db:"nullifier"`
	Type      string `db:"type"`
}

type EventsQ interface {
	New() EventsQ
	Insert(...Event) error
	Update(status EventStatus, meta json.RawMessage, points *int64) (*Event, error)
	Delete() (rowsAffected int64, err error)
	Transaction(func() error) error

	Page(*pgdb.OffsetPageParams) EventsQ
	Select() ([]Event, error)
	Get() (*Event, error)
	// Count returns the total number of events that match filters. Page is not
	// applied to this method.
	Count() (int, error)
	// SelectReopenable returns events matching criteria: there are no open or
	// fulfilled events of this type for a specific user.
	SelectReopenable() ([]ReopenableEvent, error)
	// SelectAbsentTypes returns events matching criteria: there are no events of
	// this type for a specific user. Filters are not applied to this selection.
	SelectAbsentTypes(allTypes ...string) ([]ReopenableEvent, error)

	FilterByID(...string) EventsQ
	FilterByNullifier(string) EventsQ
	FilterByStatus(...EventStatus) EventsQ
	FilterByType(...string) EventsQ
	FilterByNotType(types ...string) EventsQ
	FilterByUpdatedAtBefore(int64) EventsQ
	FilterByExternalID(string) EventsQ
	FilterInactiveNotClaimed(types ...string) EventsQ
}
