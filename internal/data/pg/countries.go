package pg

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
		updater:  squirrel.Update(countriesTable),
	}
}

func (q *countries) New() data.CountriesQ {
	return NewCountries(q.db)
}

func (q *countries) Insert(countries ...data.Country) error {
	if len(countries) == 0 {
		return nil
	}

	stmt := squirrel.Insert(countriesTable).
		Columns("code", "reserve_limit", "reserved", "withdrawn", "reserve_allowed", "withdrawal_allowed")
	for _, c := range countries {
		stmt = stmt.Values(c.Code, c.ReserveLimit, c.Reserved, c.Withdrawn, c.ReserveAllowed, c.WithdrawalAllowed)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert countries [%+v]: %w", countries, err)
	}

	return nil
}

func (q *countries) Update(fields map[string]any) error {
	if err := q.db.Exec(q.updater.SetMap(fields)); err != nil {
		return fmt.Errorf("update countries [%v]: %w", fields, err)
	}

	return nil
}

func (q *countries) UpdateMany(countries []data.Country) error {
	if len(countries) == 0 {
		return nil
	}

	values := make([]string, 0, len(countries))
	for _, v := range countries {
		values = append(values, fmt.Sprintf("('%s', %d, %t, %t)", v.Code, v.ReserveLimit, v.ReserveAllowed, v.WithdrawalAllowed))
	}

	stmt := q.updater.SetMap(map[string]interface{}{
		"reserve_limit":      squirrel.Expr("vl.reserve_limit"),
		"reserve_allowed":    squirrel.Expr("vl.reserve_allowed"),
		"withdrawal_allowed": squirrel.Expr("vl.withdrawal_allowed"),
	}).
		From(fmt.Sprintf("(VALUES %s) AS vl (code, reserve_limit, reserve_allowed, withdrawal_allowed)", strings.Join(values, ","))).
		Where(fmt.Sprintf("%s.code = vl.code", countriesTable))

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
