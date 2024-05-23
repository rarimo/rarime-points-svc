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
	return NewBalances(q.db)
}

func (q *balances) Insert(bal data.Balance) error {
	stmt := squirrel.Insert(balancesTable).SetMap(map[string]interface{}{
		"nullifier":        bal.Nullifier,
		"amount":           bal.Amount,
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

func (q *balances) SetPassport(hash string, exp time.Time, isWithdrawalAllowed bool) error {
	stmt := q.updater.
		Set("passport_hash", hash).
		Set("passport_expires", exp).
		Set("is_withdrawal_allowed", isWithdrawalAllowed)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("set passport hash and expires, and isWithdrawalAllowed: %w", err)
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

func (q *balances) Page(page *pgdb.OffsetPageParams) data.BalancesQ {
	q.selector = page.ApplyTo(q.selector, "amount", "updated_at")
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

func (q *balances) GetWithRank(nullifier string) (*data.Balance, error) {
	var res data.Balance
	stmt := fmt.Sprintf(`
	SELECT b1.*, COALESCE(b2.rank, 0) AS rank FROM %s AS b1 
	LEFT JOIN (SELECT nullifier, ROW_NUMBER() OVER (ORDER BY amount DESC, updated_at DESC) AS rank FROM %s WHERE referred_by IS NOT NULL) AS b2 
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

func (q *balances) FilterByPassportHash(passportHash string) data.BalancesQ {
	return q.applyCondition(squirrel.Eq{"passport_hash": passportHash})
}

func (q *balances) FilterDisabled() data.BalancesQ {
	return q.applyCondition(squirrel.NotEq{"referred_by": nil})
}

func (q *balances) applyCondition(cond squirrel.Sqlizer) data.BalancesQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	return q
}
