package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const balancesTable = "balances"

type balances struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
	rank     squirrel.SelectBuilder
}

func NewBalances(db *pgdb.DB) data.BalancesQ {
	return &balances{
		db:       db,
		selector: squirrel.Select("*").From(balancesTable),
		updater:  squirrel.Update(balancesTable),
		rank:     squirrel.Select("*, ROW_NUMBER() OVER (ORDER BY amount DESC, updated_at ASC) AS rank").From(balancesTable),
	}
}

func (q *balances) New() data.BalancesQ {
	return NewBalances(q.db)
}

func (q *balances) Insert(bal data.Balance) error {
	stmt := squirrel.Insert(balancesTable).SetMap(map[string]interface{}{
		"nullifier":   bal.Nullifier,
		"amount":      bal.Amount,
		"referred_by": bal.ReferredBy,
		"level":       bal.Level,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert balance %+v: %w", bal, err)
	}

	return nil
}

func (q *balances) UpdateAmountBy(points int64) error {
	stmt := q.updater.Set("amount", squirrel.Expr("amount + ?", points))

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("update amount by %d points: %w", points, err)
	}

	return nil
}

func (q *balances) SetReferredBy(referralCode string) error {
	stmt := q.updater.
		Set("referred_by", referralCode)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("set referred_by: %w", err)
	}

	return nil
}

func (q *balances) SetLevel(level int) error {
	stmt := q.updater.
		Set("level", level)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("set level: %w", err)
	}

	return nil
}

// ApplyRankedPage is similar to the ApplyTo method for a page,
// but the sorting values ​​are hardcoded because the fields must
// be sorted in opposite directions
func ApplyRankedPage(page *pgdb.OffsetPageParams, sql squirrel.SelectBuilder) squirrel.SelectBuilder {
	if page.Limit == 0 {
		page.Limit = 15
	}
	if page.Order == "" {
		page.Order = pgdb.OrderTypeDesc
	}

	offset := page.Limit * page.PageNumber

	sql = sql.Limit(page.Limit).Offset(offset)

	switch page.Order {
	case pgdb.OrderTypeAsc:
		sql = sql.OrderBy("amount asc")
		sql = sql.OrderBy("updated_at desc")
	case pgdb.OrderTypeDesc:
		sql = sql.OrderBy("amount desc")
		sql = sql.OrderBy("updated_at asc")
	default:
		panic(fmt.Errorf("unexpected order type: %v", page.Order))
	}

	return sql
}

func (q *balances) Page(page *pgdb.OffsetPageParams) data.BalancesQ {
	q.selector = ApplyRankedPage(page, q.selector)
	return q
}

func (q *balances) Select() ([]data.Balance, error) {
	var res []data.Balance

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select balances: %w", err)
	}

	return res, nil
}

func (q *balances) SelectWithRank() ([]data.Balance, error) {
	var res []data.Balance

	if err := q.db.Select(&res, q.rank); err != nil {
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

func (q *balances) GetWithRank(nullifier string) (*data.Balance, error) {
	var res data.Balance
	stmt := fmt.Sprintf(`
	SELECT b1.*, COALESCE(b2.rank, 0) AS rank FROM %s AS b1 
	LEFT JOIN (SELECT nullifier, ROW_NUMBER() OVER (ORDER BY amount DESC, updated_at ASC) AS rank FROM %s WHERE referred_by IS NOT NULL) AS b2 
	ON b1.nullifier = b2.nullifier
	WHERE b1.nullifier = ?
	`, balancesTable, balancesTable)

	if err := q.db.GetRaw(&res, stmt, nullifier); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get balance with rank: %w", err)
	}

	return &res, nil
}

func (q *balances) FilterByNullifier(nullifier string) data.BalancesQ {
	return q.applyCondition(squirrel.Eq{"nullifier": nullifier})
}

func (q *balances) FilterDisabled() data.BalancesQ {
	return q.applyCondition(squirrel.NotEq{"referred_by": nil})
}

func (q *balances) applyCondition(cond squirrel.Sqlizer) data.BalancesQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.rank = q.rank.Where(cond)
	return q
}
