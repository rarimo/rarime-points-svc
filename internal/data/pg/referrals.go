package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const referralsTable = "referrals"

type referrals struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
	consumer squirrel.UpdateBuilder
	counter  squirrel.SelectBuilder
}

func NewReferrals(db *pgdb.DB) data.ReferralsQ {
	return &referrals{
		db:       db,
		selector: squirrel.Select("id", referralsTable+".nullifier AS nullifier", "usage_left").From(referralsTable),
		updater:  squirrel.Update(referralsTable),
		consumer: squirrel.Update(referralsTable).Set("usage_left", squirrel.Expr("usage_left - 1")),
		counter:  squirrel.Select("COUNT(*) as count").From(referralsTable),
	}
}

func (q *referrals) New() data.ReferralsQ {
	return NewReferrals(q.db)
}

func (q *referrals) Insert(referrals ...data.Referral) error {
	if len(referrals) == 0 {
		return nil
	}

	stmt := squirrel.Insert(referralsTable).Columns("id", "nullifier", "usage_left")
	for _, ref := range referrals {
		stmt = stmt.Values(ref.ID, ref.Nullifier, ref.UsageLeft)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert referrals [%+v]: %w", referrals, err)
	}

	return nil
}

func (q *referrals) Update(usageLeft int) (*data.Referral, error) {
	var res data.Referral

	if err := q.db.Get(&res, q.updater.Set("usage_left", usageLeft).Suffix("RETURNING *")); err != nil {
		return nil, fmt.Errorf("update referral: %w", err)
	}

	return &res, nil
}

func (q *referrals) Consume(ids ...string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var res struct {
		IDs []string `db:"id"`
	}

	stmt := q.consumer.Where(squirrel.Eq{"id": ids}).Suffix("Returning id")

	if err := q.db.Exec(stmt); err != nil {
		return nil, fmt.Errorf("consume referrals [%v]: %w", ids, err)
	}

	return res.IDs, nil
}

func (q *referrals) ConsumeFirst(nullifier string, count uint64) error {
	stmt := fmt.Sprintf(`
		UPDATE %s SET is_consumed = true WHERE id IN (
			SELECT id FROM %s
			WHERE nullifier = ? AND is_consumed = false
			LIMIT %d
		);
	`, referralsTable, referralsTable, count)

	if err := q.db.ExecRaw(stmt, nullifier); err != nil {
		return fmt.Errorf("consume first %d referrals: %w", count, err)
	}

	return nil
}

func (q *referrals) Select() ([]data.Referral, error) {
	var res []data.Referral

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select referrals: %w", err)
	}

	return res, nil
}

func (q *referrals) Get(id string) (*data.Referral, error) {
	var res data.Referral

	if err := q.db.Get(&res, q.selector.Where(squirrel.Eq{"id": id})); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get referral by ID: %w", err)
	}

	return &res, nil
}

func (q *referrals) Count() (uint64, error) {
	var res struct {
		Count uint64 `db:"count"`
	}

	if err := q.db.Get(&res, q.counter); err != nil {
		return 0, fmt.Errorf("count referrals: %w", err)
	}

	return res.Count, nil
}

func (q *referrals) WithStatus() data.ReferralsQ {
	var (
		joinReferrer = fmt.Sprintf("JOIN %s rr ON %s.nullifier = rr.nullifier", balancesTable, referralsTable)
		joinReferee  = fmt.Sprintf("LEFT JOIN %s re ON %s.id = re.referred_by", balancesTable, referralsTable)
		joinCountry  = fmt.Sprintf("LEFT JOIN %s c ON rr.country = c.code", countriesTable)

		status = fmt.Sprintf(`CASE
			WHEN usage_left > 0 THEN '%s'
			WHEN rr.country IS NOT NULL AND NOT c.reserve_allowed AND NOT c.withdrawal_allowed THEN '%s'
			WHEN rr.country IS NOT NULL AND (c.reserved >= c.reserve_limit OR NOT c.reserve_allowed) THEN '%s'
			WHEN rr.country IS NULL AND re.country IS NOT NULL THEN '%s'
			WHEN rr.country IS NOT NULL AND re.country IS NOT NULL THEN '%s'
			ELSE '%s'
		END AS status`,
			data.StatusActive, data.StatusBanned, data.StatusLimited,
			data.StatusAwaiting, data.StatusRewarded, data.StatusConsumed,
		)
	)

	q.selector = q.selector.Column(status).
		JoinClause(joinReferrer).
		JoinClause(joinReferee).
		JoinClause(joinCountry)

	return q
}

func (q *referrals) FilterByID(code string) data.ReferralsQ {
	return q.applyCondition(squirrel.Eq{fmt.Sprintf("%s.ID", referralsTable): code})
}

func (q *referrals) FilterByNullifier(nullifier string) data.ReferralsQ {
	return q.applyCondition(squirrel.Eq{fmt.Sprintf("%s.nullifier", referralsTable): nullifier})
}

func (q *referrals) FilterConsumed() data.ReferralsQ {
	return q.applyCondition(squirrel.Gt{"usage_left": 0})
}

func (q *referrals) applyCondition(cond squirrel.Sqlizer) data.ReferralsQ {
	q.selector = q.selector.Where(cond)
	q.consumer = q.consumer.Where(cond)
	q.updater = q.updater.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
