package pg

import (
	"fmt"
	"time"

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
		Columns("balance_id", "type", "status", "created_at", "meta", "points_amount")
	for _, event := range events {
		stmt = stmt.Values(event.BalanceID, event.Type, event.Status, event.CreatedAt, event.Meta, event.PointsAmount)
	}

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert events [%+v]: %w", events, err)
	}

	return nil
}

func (q *events) Update(event data.Event) error {
	umap := map[string]any{
		"status":        event.Status,
		"meta":          event.Meta,
		"points_amount": event.PointsAmount,
		"updated_at":    time.Now().UTC(),
	}

	stmt := squirrel.Update(eventsTable).SetMap(umap).Where(squirrel.Eq{"id": event.ID})
	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("update event with map %+v: %w", umap, err)
	}

	return nil
}

func (q *events) Page(page *pgdb.CursorPageParams) data.EventsQ {
	q.selector = page.ApplyTo(q.selector, "id")
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

func (q *events) FilterByBalanceID(ids ...string) data.EventsQ {
	if len(ids) == 0 {
		return q
	}
	q.selector = q.selector.Where(squirrel.Eq{"balance_id": ids})
	q.counter = q.counter.Where(squirrel.Eq{"balance_id": ids})
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
