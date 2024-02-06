package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Balance struct {
	DID       string `db:"did"`
	Amount    int32  `db:"amount"`
	CreatedAt int32  `db:"created_at"`
	UpdatedAt int32  `db:"updated_at"`
	Rank      *int   `db:"rank"`
}

type BalancesQ interface {
	New() BalancesQ
	Insert(did string) error
	UpdateAmountBy(points int32) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)
	WithRank() BalancesQ

	FilterByDID(string) BalancesQ
}
