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
	ID           string        `db:"id"`
	UserDID      string        `db:"user_did"`
	Type         string        `db:"type"`
	Status       EventStatus   `db:"status"`
	CreatedAt    int32         `db:"created_at"`
	UpdatedAt    int32         `db:"updated_at"`
	Meta         Jsonb         `db:"meta"`
	PointsAmount sql.NullInt32 `db:"points_amount"`
}

// ReopenableEvent is a pair that is sufficient to build a new open event with a specific type for a user
type ReopenableEvent struct {
	UserDID string `db:"user_did"`
	Type    string `db:"type"`
}

type EventsQ interface {
	New() EventsQ
	Insert(...Event) error
	Update(status EventStatus, meta json.RawMessage, points *int32) (*Event, error)
	Transaction(func() error) error

	Page(*pgdb.CursorPageParams) EventsQ
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

	FilterByID(string) EventsQ
	FilterByUserDID(string) EventsQ
	FilterByStatus(...EventStatus) EventsQ
	FilterByType(...string) EventsQ
	FilterByUpdatedAtBefore(int64) EventsQ
}
