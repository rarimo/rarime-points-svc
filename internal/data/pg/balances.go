package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
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

func (q *balances) Insert(did string) error {
	stmt := squirrel.Insert(balancesTable).Columns("did").Values(did)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert balance for did %s: %w", did, err)
	}

	return nil
}

func (q *balances) UpdateAmountBy(points uint64) error {
	stmt := q.updater.Set("amount", squirrel.Expr("amount + ?", points))

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("update amount by %d points: %w", points, err)
	}

	return nil
}

func (q *balances) SetPassport(hash string, exp time.Time) error {
	stmt := q.updater.
		Set("passport_hash", hash).
		Set("passport_expires", exp)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("set passport hash and expires: %w", err)
	}

	return nil
}

func (q *balances) Page(page *pgdb.OffsetPageParams) data.BalancesQ {
	q.selector = page.ApplyTo(q.selector, "amount")
	return q
}

func (q *balances) Select() ([]data.Balance, error) {
	var res []data.Balance

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select balances: %w", err)
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

func (q *balances) WithRank() data.BalancesQ {
	q.selector = q.selector.Column("RANK() OVER (ORDER BY amount DESC, updated_at ASC) AS rank")
	return q
}

func (q *balances) FilterByDID(did string) data.BalancesQ {
	q.selector = q.selector.Where(squirrel.Eq{"did": did})
	q.updater = q.updater.Where(squirrel.Eq{"did": did})
	return q
}
