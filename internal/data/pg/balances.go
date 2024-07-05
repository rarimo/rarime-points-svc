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
	counter  squirrel.SelectBuilder
}

func NewBalances(db *pgdb.DB) data.BalancesQ {
	return &balances{
		db:       db,
		selector: squirrel.Select("*").From(balancesTable),
		updater:  squirrel.Update(balancesTable),
		rank:     squirrel.Select("*, ROW_NUMBER() OVER (ORDER BY amount DESC, updated_at ASC) AS rank").From(balancesTable),
		counter:  squirrel.Select("COUNT(*) as count").From(balancesTable),
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
		"country":     bal.Country,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert balance %+v: %w", bal, err)
	}

	return nil
}

func (q *balances) Update(fields map[string]any) error {
	if err := q.db.Exec(q.updater.SetMap(fields)); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}

// applyRankedPage is similar to the pgdb.OffsetParams.ApplyTo method,
// but the sorting values are hardcoded because the fields must
// be sorted in opposite directions
func applyRankedPage(page *pgdb.OffsetPageParams, sql squirrel.SelectBuilder) squirrel.SelectBuilder {
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
	q.selector = applyRankedPage(page, q.selector)
	q.rank = applyRankedPage(page, q.rank)
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

func (q *balances) Count() (int64, error) {
	res := struct {
		Count int64 `db:"count"`
	}{}

	if err := q.db.Get(&res, q.counter); err != nil {
		return 0, fmt.Errorf("get balance: %w", err)
	}

	return res.Count, nil
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

func (q *balances) WithoutPassportEvent() ([]data.WithoutPassportEventBalance, error) {
	var res []data.WithoutPassportEventBalance
	stmt := fmt.Sprintf(`
	SELECT b.*, e.id AS event_id, e.status AS event_status
		FROM %s AS b INNER JOIN %s AS e
		ON b.nullifier = e.nullifier AND e.type='passport_scan' 
		WHERE e.status NOT IN ('claimed') 
		AND b.referred_by IS NOT NULL
		AND b.country IS NOT NULL
	`, balancesTable, eventsTable)

	if err := q.db.SelectRaw(&res, stmt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select balances without passport events: %w", err)
	}

	return res, nil
}

func (q *balances) WithoutReferralEvent() ([]data.ReferredReferrer, error) {
	var res []data.ReferredReferrer
	stmt := fmt.Sprintf(`
	SELECT b.nullifier AS referred, r.nullifier AS referrer 
		FROM %s AS b INNER JOIN %s AS r 
		ON r.id = b.referred_by 
		WHERE b.nullifier NOT IN 
			(SELECT b.nullifier 
				FROM %s AS b INNER JOIN %s AS e 
				ON e.meta->>'nullifier' = b.nullifier) 
		AND b.referred_by IS NOT NULL 
		AND b.country IS NOT NULL
	`, balancesTable, referralsTable, balancesTable, eventsTable)

	if err := q.db.SelectRaw(&res, stmt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select balances without referred events: %w", err)
	}

	return res, nil

}

func (q *balances) FilterByNullifier(nullifier ...string) data.BalancesQ {
	return q.applyCondition(squirrel.Eq{"nullifier": nullifier})
}

func (q *balances) FilterDisabled() data.BalancesQ {
	return q.applyCondition(squirrel.NotEq{"referred_by": nil})
}

func (q *balances) FilterByAnonymousID(id string) data.BalancesQ {
	return q.applyCondition(squirrel.Eq{"anonymous_id": id})
}

func (q *balances) applyCondition(cond squirrel.Sqlizer) data.BalancesQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.rank = q.rank.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
