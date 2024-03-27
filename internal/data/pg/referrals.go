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
	counter  squirrel.SelectBuilder
}

func NewReferrals(db *pgdb.DB) data.ReferralsQ {
	return &referrals{
		db:       db,
		selector: squirrel.Select("*").From(referralsTable),
		updater:  squirrel.Update(referralsTable).Set("is_consumed", true),
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

	stmt := squirrel.Insert(referralsTable).Columns("id", "user_did")
	for _, ref := range referrals {
		stmt = stmt.Values(ref.ID, ref.UserDID)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert referrals [%+v]: %w", referrals, err)
	}

	return nil
}

func (q *referrals) Consume(ids ...string) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var res struct {
		IDs []string `db:"id"`
	}

	stmt := q.updater.Where(squirrel.Eq{"id": ids}).Suffix("Returning id")

	if err := q.db.Exec(stmt); err != nil {
		return nil, fmt.Errorf("consume referrals [%v]: %w", ids, err)
	}

	return res.IDs, nil
}

func (q *referrals) ConsumeFirst(did string, count uint64) error {
	stmt := fmt.Sprintf(`
		UPDATE %s SET is_consumed = true WHERE id IN (
			SELECT id FROM %s
			WHERE user_did = ? AND is_consumed = false
			LIMIT %d
		);
	`, referralsTable, referralsTable, count)

	if err := q.db.ExecRaw(stmt, did); err != nil {
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

func (q *referrals) FilterByUserDID(did string) data.ReferralsQ {
	return q.applyCondition(squirrel.Eq{"user_did": did})
}

func (q *referrals) FilterByIsConsumed(isConsumed bool) data.ReferralsQ {
	return q.applyCondition(squirrel.Eq{"is_consumed": isConsumed})
}

func (q *referrals) applyCondition(cond squirrel.Sqlizer) data.ReferralsQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
