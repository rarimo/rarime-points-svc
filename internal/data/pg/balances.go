package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const balancesTable = "balances"

type balances struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewBalances(db *pgdb.DB) data.BalancesQ {
	return &balances{
		db:       db,
		selector: squirrel.Select("*").From(balancesTable),
		updater:  squirrel.Update(balancesTable),
	}
}

func (q *balances) New() data.BalancesQ {
	return NewBalances(q.db.Clone())
}

func (q *balances) Insert(balance data.Balance) error {
	stmt := squirrel.Insert(balancesTable).SetMap(map[string]interface{}{
		"id":     balance.ID,
		"did":    balance.DID,
		"amount": balance.Amount,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert balance %+v: %w", balance, err)
	}

	return nil
}

func (q *balances) UpdateAmount(amount int) error {
	stmt := q.updater.Set("amount", amount)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("update balance amount to %d: %w", amount, err)
	}

	return nil
}

func (q *balances) SelectLeaders(count int) ([]data.Balance, error) {
	var res []data.Balance

	stmt := squirrel.Select("*").
		From(balancesTable).
		OrderBy("amount DESC, updated_at ASC").
		Limit(uint64(count))

	if err := q.db.Select(&res, stmt); err != nil {
		return nil, fmt.Errorf("select leaders: %w", err)
	}

	return res, nil
}

func (q *balances) Get() (*data.Balance, error) {
	var res data.Balance

	if err := q.db.Get(&res, q.selector); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return &res, nil
}

func (q *balances) FilterByID(id string) data.BalancesQ {
	q.selector = q.selector.Where(squirrel.Eq{"id": id})
	q.updater = q.updater.Where(squirrel.Eq{"id": id})
	return q
}

func (q *balances) FilterByUserDID(did string) data.BalancesQ {
	q.selector = q.selector.Where(squirrel.Eq{"did": did})
	q.updater = q.updater.Where(squirrel.Eq{"did": did})
	return q
}
