package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const faceEventBalancesTable = "face_event_balances"

type faceEventBalances struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewFaceEventBalances(db *pgdb.DB) data.FaceEventBalanceQ {
	return &faceEventBalances{
		db:       db,
		selector: squirrel.Select("*").From(faceEventBalancesTable),
		updater:  squirrel.Update(faceEventBalancesTable),
	}
}

func (q *faceEventBalances) New() data.FaceEventBalanceQ {
	return NewFaceEventBalances(q.db)
}

func (q *faceEventBalances) Insert(bal data.FaceEventBalance) error {
	stmt := squirrel.Insert(faceEventBalancesTable).SetMap(map[string]interface{}{
		"nullifier": bal.Nullifier,
		"amount":    bal.Amount,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert balance %+v: %w", bal, err)
	}

	return nil
}

func (q *faceEventBalances) Update(fields map[string]any) error {
	if err := q.db.Exec(q.updater.SetMap(fields)); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}

func (q *faceEventBalances) Get() (*data.FaceEventBalance, error) {
	var res data.FaceEventBalance

	if err := q.db.Get(&res, q.selector); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get balance: %w", err)
	}

	return &res, nil
}

func (q *faceEventBalances) Transaction(f func() error) error {
	return q.db.Transaction(f)
}

func (q *faceEventBalances) applyCondition(cond squirrel.Sqlizer) data.FaceEventBalanceQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	return q
}

func (q *faceEventBalances) FilterByNullifier(nullifier ...string) data.FaceEventBalanceQ {
	return q.applyCondition(squirrel.Eq{"nullifier": nullifier})
}
