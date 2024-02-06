package pg

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const withdrawalsTable = "withdrawals"

type withdrawals struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewWithdrawals(db *pgdb.DB) data.WithdrawalsQ {
	return &withdrawals{
		db:       db,
		selector: squirrel.Select("*").From(withdrawalsTable),
	}
}

func (q *withdrawals) New() data.WithdrawalsQ {
	return NewWithdrawals(q.db.Clone())
}

func (q *withdrawals) Insert(w data.Withdrawal) error {
	stmt := squirrel.Insert(withdrawalsTable).SetMap(map[string]interface{}{
		"user_did": w.UserDID,
		"amount":   w.Amount,
		"address":  w.Address,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert withdrawal [%+v]: %w", w, err)
	}

	return nil
}

func (q *withdrawals) Page(page *pgdb.CursorPageParams) data.WithdrawalsQ {
	q.selector = page.ApplyTo(q.selector, "created_at")
	return q
}

func (q *withdrawals) Select() ([]data.Withdrawal, error) {
	var res []data.Withdrawal

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select withdrawals: %w", err)
	}

	return res, nil
}

func (q *withdrawals) FilterByUserDID(did string) data.WithdrawalsQ {
	q.selector = q.selector.Where(squirrel.Eq{"user_did": did})
	return q
}
