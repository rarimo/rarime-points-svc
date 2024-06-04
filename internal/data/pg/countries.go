package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const countriesTable = "countries"

type countries struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewCountries(db *pgdb.DB) data.CountriesQ {
	return &countries{
		db:       db,
		selector: squirrel.Select("*").From(countriesTable),
		updater:  squirrel.Update(countriesTable).Set("usage_left", squirrel.Expr("usage_left - 1")),
	}
}

func (q *countries) New() data.CountriesQ {
	return NewCountries(q.db)
}

func (q *countries) Insert(countries ...data.Country) error {
	if len(countries) == 0 {
		return nil
	}

	stmt := squirrel.Insert(countriesTable).Columns("code", "reserve_limit", "reserved", "withdrawn", "is_disabled")
	for _, c := range countries {
		stmt = stmt.Values(c.Code, c.ReserveLimit, c.Reserved, c.Withdrawn, c.IsDisabled)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert countries [%+v]: %w", countries, err)
	}

	return nil
}

func (q *countries) Update(limit, addReserved, addWithdrawn *int64, isDisabled *bool) error {
	stmt := q.updater

	if limit != nil {
		stmt = stmt.Set("usage_left", *limit)
	}
	if addReserved != nil {
		stmt = stmt.Set("reserved", squirrel.Expr("reserved + ?", *addReserved))
	}
	if addWithdrawn != nil {
		stmt = stmt.Set("withdrawn", squirrel.Expr("withdrawn + ?", *addWithdrawn))
	}
	if isDisabled != nil {
		stmt = stmt.Set("is_disabled", *isDisabled)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("update countries: %w", err)
	}

	return nil
}

func (q *countries) Select() ([]data.Country, error) {
	var res []data.Country

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select countries: %w", err)
	}

	return res, nil
}

func (q *countries) Get() (*data.Country, error) {
	var res data.Country

	if err := q.db.Get(&res, q.selector); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get country: %w", err)
	}

	return &res, nil
}

func (q *countries) FilterByCodes(codes ...string) data.CountriesQ {
	return q.applyCondition(squirrel.Eq{"code": codes})
}

func (q *countries) applyCondition(cond squirrel.Sqlizer) data.CountriesQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	return q
}
