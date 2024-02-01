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

type EventsQ interface {
	New() EventsQ
	Insert(...Event) error
	Update(status EventStatus, meta []json.RawMessage, points *int32) (*Event, error)

	Page(*pgdb.CursorPageParams) EventsQ
	Select() ([]Event, error)
	Get() (*Event, error)
	Count() (int, error)

	FilterByID(string) EventsQ
	FilterByUserDID(string) EventsQ
	FilterByStatus(...EventStatus) EventsQ
	FilterByType(...string) EventsQ
}

type BalancesQ interface {
	New() BalancesQ
	Insert(did string) error
	UpdateAmountBy(points int32) error
	SetAddress(string) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)
	WithRank() BalancesQ

	FilterByDID(string) BalancesQ
}

type Event struct {
	ID           string          `db:"id"`
	UserDID      string          `db:"user_did"`
	Type         string          `db:"type"`
	Status       EventStatus     `db:"status"`
	CreatedAt    int32           `db:"created_at"`
	UpdatedAt    int32           `db:"updated_at"`
	Meta         json.RawMessage `db:"meta"`
	PointsAmount sql.NullInt32   `db:"points_amount"`
}

type Balance struct {
	DID       string         `db:"did"`
	Amount    int32          `db:"amount"`
	CreatedAt int32          `db:"created_at"`
	UpdatedAt int32          `db:"updated_at"`
	Address   sql.NullString `db:"address"`
	Rank      *int           `db:"rank"`
}
