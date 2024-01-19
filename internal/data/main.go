package data

import (
	"database/sql"
	"time"
)

type EventsQ interface {
	New() EventsQ
	Insert(Event) error
	UpdateStatus(string) error
	Select() ([]Event, error)

	FilterByID(string) EventsQ
	FilterByBalanceID(...string) EventsQ
	FilterByStatus(...string) EventsQ
}

type BalancesQ interface {
	New() BalancesQ
	Insert(Balance) error
	UpdateAmount(int) error
	Get() (*Balance, error)

	FilterByID(string) BalancesQ
	FilterByUserDID(string) BalancesQ
}

type Event struct {
	ID        string         `db:"id"`
	TypeID    string         `db:"type_id"`
	BalanceID string         `db:"balance_id"`
	Status    string         `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
	Meta      sql.NullString `db:"meta"`
}

type Balance struct {
	ID        string    `db:"id"`
	DID       string    `db:"did"`
	Amount    int       `db:"amount"`
	UpdatedAt time.Time `db:"updated_at"`
}
