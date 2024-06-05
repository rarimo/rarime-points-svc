package data

import (
	"database/sql"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	ColAmount     = "amount"
	ColReferredBy = "referred_by"
	ColLevel      = "level"
	ColCountry    = "country"
)

type Balance struct {
	Nullifier  string         `db:"nullifier"`
	Amount     int64          `db:"amount"`
	CreatedAt  int32          `db:"created_at"`
	UpdatedAt  int32          `db:"updated_at"`
	ReferredBy sql.NullString `db:"referred_by"`
	Rank       *int           `db:"rank"`
	Level      int            `db:"level"`
	Country    *string        `db:"country"`
}

type BalancesQ interface {
	New() BalancesQ
	Insert(Balance) error
	Update(map[string]any) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)
	// GetWithRank returns balance with rank, filtered by nullifier. No other filters can be applied.
	GetWithRank(nullifier string) (*Balance, error)
	SelectWithRank() ([]Balance, error)

	FilterByNullifier(string) BalancesQ
	FilterDisabled() BalancesQ
}
