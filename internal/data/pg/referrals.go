package pg

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const referralsTable = "referrals"

type referrals struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewReferrals(db *pgdb.DB) data.ReferralsQ {
	return &referrals{
		db:       db,
		selector: squirrel.Select("*").From(referralsTable),
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

func (q *referrals) Deactivate(id string) error {
	stmt := squirrel.Update(referralsTable).
		Set("is_consumed", true).
		Where(squirrel.Eq{"id": id})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("deactivate referral [id=%s]: %w", id, err)
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
		return nil, fmt.Errorf("get referral by ID: %w", err)
	}

	return &res, nil
}

func (q *referrals) FilterByUserDID(did string) data.ReferralsQ {
	q.selector = q.selector.Where(squirrel.Eq{"user_did": did})
	return q
}

func (q *referrals) FilterByIsConsumed(isConsumed bool) data.ReferralsQ {
	q.selector = q.selector.Where(squirrel.Eq{"is_consumed": isConsumed})
	return q
}
