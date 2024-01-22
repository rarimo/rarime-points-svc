package pg

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const eventsTable = "events"

type events struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewEvents(db *pgdb.DB) data.EventsQ {
	return &events{
		db:       db,
		selector: squirrel.Select("*").From(eventsTable),
		updater:  squirrel.Update(eventsTable),
	}
}

func (q *events) New() data.EventsQ {
	return NewEvents(q.db.Clone())
}

func (q *events) Insert(event data.Event) error {
	stmt := squirrel.Insert(eventsTable).SetMap(map[string]interface{}{
		"id":         event.ID,
		"type_id":    event.TypeID,
		"balance_id": event.BalanceID,
		"status":     event.Status,
		"created_at": event.CreatedAt,
		"meta":       event.Meta,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert event %+v: %w", event, err)
	}

	return nil
}

func (q *events) UpdateStatus(status data.EventStatus) error {
	stmt := q.updater.Set("status", status)

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("update event status to %s: %w", status, err)
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
	q.updater = q.updater.Where(squirrel.Eq{"id": id})
	return q
}

func (q *events) FilterByBalanceID(ids ...string) data.EventsQ {
	q.selector = q.selector.Where(squirrel.Eq{"balance_id": ids})
	q.updater = q.updater.Where(squirrel.Eq{"balance_id": ids})
	return q
}

func (q *events) FilterByStatus(statuses ...data.EventStatus) data.EventsQ {
	q.selector = q.selector.Where(squirrel.Eq{"status": statuses})
	q.updater = q.updater.Where(squirrel.Eq{"status": statuses})
	return q
}
