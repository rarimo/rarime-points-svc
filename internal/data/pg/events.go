package pg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const eventsTable = "events"

type events struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
	counter  squirrel.SelectBuilder
}

func NewEvents(db *pgdb.DB) data.EventsQ {
	return &events{
		db:       db,
		selector: squirrel.Select("*").From(eventsTable),
		updater:  squirrel.Update(eventsTable),
		counter:  squirrel.Select("count(id) AS count").From(eventsTable),
	}
}

func (q *events) New() data.EventsQ {
	return NewEvents(q.db.Clone())
}

func (q *events) Insert(events ...data.Event) error {
	if len(events) == 0 {
		return nil
	}

	stmt := squirrel.Insert(eventsTable).
		Columns("user_did", "type", "status", "meta", "points_amount")
	for _, event := range events {
		var meta any
		if len(event.Meta) != 0 {
			meta = event.Meta
		}
		stmt = stmt.Values(event.UserDID, event.Type, event.Status, meta, event.PointsAmount)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert events [%+v]: %w", events, err)
	}

	return nil
}

func (q *events) Update(status data.EventStatus, meta json.RawMessage, points *int32) (*data.Event, error) {
	umap := map[string]any{
		"status": status,
	}
	if points != nil {
		umap["points_amount"] = sql.NullInt32{Int32: *points, Valid: true}
	}
	if len(meta) != 0 {
		umap["meta"] = meta
	}

	var res data.Event
	stmt := q.updater.SetMap(umap).Suffix("RETURNING *")

	if err := q.db.Get(&res, stmt); err != nil {
		return nil, fmt.Errorf("update event with map %+v: %w", umap, err)
	}

	return &res, nil
}

func (q *events) Page(page *pgdb.CursorPageParams) data.EventsQ {
	q.selector = page.ApplyTo(q.selector, "updated_at")
	return q
}

func (q *events) Transaction(f func() error) error {
	return q.db.Transaction(f)
}

func (q *events) Select() ([]data.Event, error) {
	var res []data.Event

	if err := q.db.Select(&res, q.selector); err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}

	return res, nil
}

func (q *events) Get() (*data.Event, error) {
	var res data.Event

	if err := q.db.Get(&res, q.selector); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get event: %w", err)
	}

	return &res, nil
}

func (q *events) Count() (int, error) {
	var res struct {
		Count int `db:"count"`
	}

	if err := q.db.Get(&res, q.counter); err != nil {
		return 0, fmt.Errorf("count events: %w", err)
	}

	return res.Count, nil
}

func (q *events) FilterByID(id string) data.EventsQ {
	return q.applyCondition(squirrel.Eq{"id": id})
}

func (q *events) FilterByUserDID(did string) data.EventsQ {
	return q.applyCondition(squirrel.Eq{"user_did": did})
}

func (q *events) FilterByStatus(statuses ...data.EventStatus) data.EventsQ {
	if len(statuses) == 0 {
		return q
	}
	return q.applyCondition(squirrel.Eq{"status": statuses})
}

func (q *events) FilterByType(types ...string) data.EventsQ {
	if len(types) == 0 {
		return q
	}
	return q.applyCondition(squirrel.Eq{"type": types})
}

func (q *events) applyCondition(cond squirrel.Eq) data.EventsQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.counter = q.counter.Where(cond)
	return q
}
