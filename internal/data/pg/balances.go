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
const balancesRankColumns = "did, MAX(amount) as amount, created_at, updated_at, referral_id, referred_by, passport_hash, passport_expires"

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
	return NewBalances(q.db)
}

func (q *balances) Insert(bal data.Balance) error {
	stmt := squirrel.Insert(balancesTable).SetMap(map[string]interface{}{
		"did":              bal.DID,
		"amount":           bal.Amount,
		"referral_id":      bal.ReferralID,
		"referred_by":      bal.ReferredBy,
		"passport_hash":    bal.PassportHash,
		"passport_expires": bal.PassportExpires,
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

func (q *balances) GetWithRank(did string) (*data.Balance, error) {
	var res data.Balance
	stmt := fmt.Sprintf(`
		SELECT * FROM (
			SELECT *, RANK() OVER (ORDER BY amount DESC, created_at ASC) AS rank FROM (
				SELECT %s FROM %s GROUP BY did
			) AS t
		) AS ranked WHERE did = ?
	`, balancesRankColumns, balancesTable)

	if err := q.db.GetRaw(&res, stmt, did); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get balance with rank: %w", err)
	}

	return &res, nil
}

func (q *balances) FilterByDID(did string) data.BalancesQ {
	return q.applyCondition(squirrel.Eq{"did": did})
}

func (q *balances) FilterByReferralID(referralID string) data.BalancesQ {
	return q.applyCondition(squirrel.Eq{"referral_id": referralID})
}

func (q *balances) applyCondition(cond squirrel.Eq) data.BalancesQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	return q
}
