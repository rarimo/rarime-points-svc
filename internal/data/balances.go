package data

import (
	"database/sql"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Balance struct {
	DID             string         `db:"did"`
	Amount          uint64         `db:"amount"`
	CreatedAt       int32          `db:"created_at"`
	UpdatedAt       int32          `db:"updated_at"`
	PassportHash    sql.NullString `db:"passport_hash"`
	PassportExpires sql.NullTime   `db:"passport_expires"`
	Rank            *int           `db:"rank"`
}

type BalancesQ interface {
	New() BalancesQ
	Insert(did string) error
	UpdateAmountBy(points uint64) error
	SetPassport(hash string, exp time.Time) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)
	WithRank() BalancesQ

	FilterByDID(string) BalancesQ
}
