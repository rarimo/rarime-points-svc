package data

import (
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Withdrawal struct {
	ID        string `db:"id"`
	UserDID   string `db:"user_did"`
	Amount    int32  `db:"amount"`
	Address   string `db:"address"`
	CreatedAt int32  `db:"created_at"`
}

type WithdrawalsQ interface {
	New() WithdrawalsQ
	Insert(Withdrawal) error

	Page(*pgdb.CursorPageParams) WithdrawalsQ
	Select() ([]Withdrawal, error)

	FilterByUserDID(string) WithdrawalsQ
}
