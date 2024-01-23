package data

import (
	"database/sql"
	"time"

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
	Update(Event) error

	Page(*pgdb.CursorPageParams) EventsQ
	Select() ([]Event, error)
	Get() (*Event, error)
	Count() (int, error)

	FilterByID(string) EventsQ
	FilterByBalanceID(...string) EventsQ
	FilterByStatus(...EventStatus) EventsQ
}

type BalancesQ interface {
	New() BalancesQ
	Insert(Balance) error
	UpdateAmount(int) error

	SelectLeaders(count int) ([]Balance, error)
	Get() (*Balance, error)
	WithRank() BalancesQ

	FilterByID(string) BalancesQ
	FilterByUserDID(string) BalancesQ
}

type Event struct {
	ID           string         `db:"id"`
	BalanceID    string         `db:"balance_id"`
	Type         string         `db:"type"`
	Status       EventStatus    `db:"status"`
	CreatedAt    time.Time      `db:"created_at"`
	Meta         sql.NullString `db:"meta"`
	PointsAmount sql.NullInt32  `db:"points_amount"`
}

type Balance struct {
	ID        string    `db:"id"`
	DID       string    `db:"did"`
	Amount    int       `db:"amount"`
	UpdatedAt time.Time `db:"updated_at"`
	Rank      *int      `db:"rank"`
}
