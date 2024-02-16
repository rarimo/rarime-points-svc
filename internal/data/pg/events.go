package pg

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const eventsTable = "events"

type events struct {
	db         *pgdb.DB
	selector   squirrel.SelectBuilder
	updater    squirrel.UpdateBuilder
	deleter    squirrel.DeleteBuilder
	counter    squirrel.SelectBuilder
	reopenable squirrel.SelectBuilder
}

func NewEvents(db *pgdb.DB) data.EventsQ {
	return &events{
		db:         db,
		selector:   squirrel.Select("*").From(eventsTable),
		updater:    squirrel.Update(eventsTable),
		deleter:    squirrel.Delete(eventsTable),
		counter:    squirrel.Select("count(id) AS count").From(eventsTable),
		reopenable: squirrel.Select("user_did", "type").Distinct().From(eventsTable + " e1"),
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

func (q *events) Update(status data.EventStatus, meta json.RawMessage, points *int64) (*data.Event, error) {
	umap := map[string]any{
		"status": status,
	}
	if len(meta) != 0 {
		umap["meta"] = meta
	}
	if points != nil {
		umap["points_amount"] = points
	}

	var res data.Event
	stmt := q.updater.SetMap(umap).Suffix("RETURNING *")

	if err := q.db.Get(&res, stmt); err != nil {
		return nil, fmt.Errorf("update event with map %+v: %w", umap, err)
	}

	return &res, nil
}

func (q *events) Delete() (int64, error) {
	res, err := q.db.ExecWithResult(q.deleter)
	if err != nil {
		return 0, fmt.Errorf("delete events: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("count rows affected: %w", err)
	}

	return rows, nil
}

func (q *events) Transaction(f func() error) error {
	return q.db.Transaction(f)
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

func (q *events) SelectReopenable() ([]data.ReopenableEvent, error) {
	subq := fmt.Sprintf(`NOT EXISTS (
	SELECT 1 FROM %s e2
    WHERE e2.user_did = e1.user_did
    AND e2.type = e1.type
    AND e2.status IN (?, ?))`, eventsTable)
	stmt := q.reopenable.Where(subq, data.EventOpen, data.EventFulfilled)

	var res []data.ReopenableEvent
	if err := q.db.Select(&res, stmt); err != nil {
		return nil, fmt.Errorf("select reopenable events: %w", err)
	}

	return res, nil
}

func (q *events) SelectAbsentTypes(allTypes ...string) ([]data.ReopenableEvent, error) {
	values := make([]string, len(allTypes))
	for i, t := range allTypes {
		values[i] = fmt.Sprintf("('%s')", t)
	}

	query := fmt.Sprintf(`
		WITH types(type) AS (
    		VALUES %s
		)
		SELECT u.user_did, t.type
		FROM (
    		SELECT DISTINCT user_did FROM %s
		) u
		CROSS JOIN types t
		LEFT JOIN %s e ON e.user_did = u.user_did AND e.type = t.type
		WHERE e.type IS NULL;
	`, strings.Join(values, ", "), eventsTable, eventsTable)

	var res []data.ReopenableEvent
	if err := q.db.SelectRaw(&res, query); err != nil {
		return nil, fmt.Errorf("select absent types for each did: %w", err)
	}

	return res, nil
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

func (q *events) FilterByUpdatedAtBefore(unix int64) data.EventsQ {
	return q.applyCondition(squirrel.Lt{"updated_at": unix})
}

func (q *events) applyCondition(cond squirrel.Sqlizer) data.EventsQ {
	q.selector = q.selector.Where(cond)
	q.updater = q.updater.Where(cond)
	q.deleter = q.deleter.Where(cond)
	q.counter = q.counter.Where(cond)
	q.reopenable = q.reopenable.Where(cond)
	return q
}
