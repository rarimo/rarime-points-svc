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
	EventReserved  EventStatus = "reserved"
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

type EventsQ interface {
	New() EventsQ
	Insert(...Event) error
	Update(status EventStatus, meta json.RawMessage, points *int32) (*Event, error)
	Reopen() (count uint, err error)
	Transaction(func() error) error

	Page(*pgdb.CursorPageParams) EventsQ
	Select() ([]Event, error)
	Get() (*Event, error)
	Count() (int, error)

	FilterByID(string) EventsQ
	FilterByUserDID(string) EventsQ
	FilterByStatus(...EventStatus) EventsQ
	FilterByType(...string) EventsQ
	FilterByUpdatedAtBefore(int64) EventsQ
}