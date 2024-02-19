package data

import (
	"database/sql"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Balance struct {
	DID             string         `db:"did"`
	Amount          int64          `db:"amount"`
	CreatedAt       int32          `db:"created_at"`
	UpdatedAt       int32          `db:"updated_at"`
	ReferralID      string         `db:"referral_id"`
	ReferredBy      sql.NullString `db:"referred_by"`
	PassportHash    sql.NullString `db:"passport_hash"`
	PassportExpires sql.NullTime   `db:"passport_expires"`
	Rank            *int           `db:"rank"`
}

type BalancesQ interface {
	New() BalancesQ
	Insert(Balance) error
	UpdateAmountBy(points int64) error
	SetPassport(hash string, exp time.Time) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)

	WithRank() BalancesQ
	FilterByDID(string) BalancesQ
	FilterByReferralID(string) BalancesQ
}
