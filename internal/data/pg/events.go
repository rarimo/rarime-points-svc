package pg

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const eventsTable = "events"

type events struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewEvents(db *pgdb.DB) data.EventsQ {
	return &events{
		db:       db,
		selector: squirrel.Select("*").From(eventsTable),
	}
}

func (q *events) New() data.EventsQ {
	return NewEvents(q.db.Clone())
}

func (q *events) Insert(event data.Event) error {
	imap := map[string]any{ // ID must be created sequentially
		"balance_id":    event.BalanceID,
		"type":          event.Type,
		"status":        event.Status,
		"created_at":    event.CreatedAt,
		"meta":          event.Meta,
		"points_amount": event.PointsAmount,
	}

	if err := q.db.Exec(squirrel.Insert(eventsTable).SetMap(imap)); err != nil {
		return fmt.Errorf("insert event %+v: %w", event, err)
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

func (q *events) FilterByID(id string) data.EventsQ {
	q.selector = q.selector.Where(squirrel.Eq{"id": id})
	return q
}

func (q *events) FilterByBalanceID(ids ...string) data.EventsQ {
	q.selector = q.selector.Where(squirrel.Eq{"balance_id": ids})
	return q
}

func (q *events) FilterByStatus(statuses ...data.EventStatus) data.EventsQ {
	q.selector = q.selector.Where(squirrel.Eq{"status": statuses})
	return q
}
