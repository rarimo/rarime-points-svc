-- +migrate Up
CREATE TABLE IF NOT EXISTS balances
(
    id         serial PRIMARY KEY,
    did        text                        not null unique,
    amount     integer                     not null default 0,
    updated_at timestamp without time zone not null default NOW()
);

CREATE TYPE event_status AS ENUM ('open', 'fulfilled', 'claimed');

CREATE TABLE IF NOT EXISTS events
(
    id            serial PRIMARY KEY,
    balance_id integer not null REFERENCES balances (id),
    type          text                        not null,
    status        event_status                not null,
    created_at    timestamp without time zone not null default NOW(),
    meta          text,
    points_amount integer
);

CREATE INDEX IF NOT EXISTS events_balance_id_index ON events using btree (balance_id);
CREATE INDEX IF NOT EXISTS events_type_index ON events using btree (type);

-- +migrate Down
DROP INDEX IF EXISTS events_type_index;
DROP INDEX IF EXISTS events_balance_id_index;
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_status;
DROP TABLE IF EXISTS balances;
