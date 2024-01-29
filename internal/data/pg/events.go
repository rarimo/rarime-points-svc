package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const eventsTable = "events"

type events struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	counter  squirrel.SelectBuilder
}

func NewEvents(db *pgdb.DB) data.EventsQ {
	return &events{
		db:       db,
		selector: squirrel.Select("*").From(eventsTable),
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
		stmt = stmt.Values(event.UserDID, event.Type, event.Status, event.Meta, event.PointsAmount)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert events [%+v]: %w", events, err)
	}

	return nil
}

func (q *events) Update(status data.EventStatus, meta []json.RawMessage, points *int32) (*data.Event, error) {
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
	stmt := squirrel.Update(eventsTable).SetMap(umap)

	if err := q.db.Get(&res, stmt); err != nil {
		return nil, fmt.Errorf("update event with map %+v: %w", umap, err)
	}

	return &res, nil
}

func (q *events) Page(page *pgdb.CursorPageParams) data.EventsQ {
	q.selector = page.ApplyTo(q.selector, "updated_at")
	return q
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
	q.selector = q.selector.Where(squirrel.Eq{"id": id})
	q.counter = q.counter.Where(squirrel.Eq{"id": id})
	return q
}

func (q *events) FilterByUserDID(did string) data.EventsQ {
	q.selector = q.selector.Where(squirrel.Eq{"user_did": did})
	q.counter = q.counter.Where(squirrel.Eq{"user_did": did})
	return q
}

func (q *events) FilterByStatus(statuses ...data.EventStatus) data.EventsQ {
	if len(statuses) == 0 {
		return q
	}
	q.selector = q.selector.Where(squirrel.Eq{"status": statuses})
	q.counter = q.counter.Where(squirrel.Eq{"status": statuses})
	return q
}

func (q *events) FilterByType(types ...string) data.EventsQ {
	if len(types) == 0 {
		return q
	}
	q.selector = q.selector.Where(squirrel.Eq{"type": types})
	q.counter = q.counter.Where(squirrel.Eq{"type": types})
	return q
}
