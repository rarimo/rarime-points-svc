package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Withdrawal struct {
	ID        string `db:"id"`
	Nullifier string `db:"nullifier"`
	Amount    int64  `db:"amount"`
	Address   string `db:"address"`
	CreatedAt int32  `db:"created_at"`
}

type WithdrawalsQ interface {
	New() WithdrawalsQ
	Insert(Withdrawal) (*Withdrawal, error)

	Page(*pgdb.CursorPageParams) WithdrawalsQ
	Select() ([]Withdrawal, error)

	FilterByNullifier(string) WithdrawalsQ
}
