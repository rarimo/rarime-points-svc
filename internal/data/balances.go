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
	ReferredBy      sql.NullString `db:"referred_by"`
	PassportHash    sql.NullString `db:"passport_hash"`
	PassportExpires sql.NullTime   `db:"passport_expires"`
	Rank            *int           `db:"rank"`
	Level           int            `db:"level"`
	LevelClaimId    *string        `db:"level_claim_id"`
}

type BalancesQ interface {
	New() BalancesQ
	Insert(Balance) error
	UpdateAmountBy(points int64) error
	SetPassport(hash string, exp time.Time) error
	SetReferredBy(referralCode string) error
	SetLevel(level int, id string) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)
	// GetWithRank returns balance with rank, filtered by DID. No other filters can be applied.
	GetWithRank(did string) (*Balance, error)

	FilterByDID(string) BalancesQ
	FilterDisabled() BalancesQ

	Transaction(func() error) error
}
