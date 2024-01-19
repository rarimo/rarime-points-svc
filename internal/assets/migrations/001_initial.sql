-- +migrate Up
DROP FUNCTION IF EXISTS trigger_set_updated_at cascade;

CREATE FUNCTION trigger_set_updated_at() RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$ BEGIN NEW.updated_at = NOW() at time zone 'utc'; RETURN NEW; END; $$;

CREATE TABLE IF NOT EXISTS balances
(
    id         serial PRIMARY KEY,
    did        text                        not null unique,
    amount     integer                     not null default 0,
    updated_at timestamp without time zone not null default NOW()
);

CREATE TRIGGER set_updated_at
    before update
    on balances
    for each row
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TYPE event_status AS ENUM ('open', 'fulfilled', 'claimed');

CREATE TABLE IF NOT EXISTS events
(
    id         serial PRIMARY KEY,
    type_id    smallint                    not null,
    balance_id integer                     null REFERENCES balances (id),
    status     event_status                not null,
    created_at timestamp without time zone not null default NOW(),
    meta       text
);

-- +migrate Down
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_status;
DROP TRIGGER IF EXISTS set_updated_at on balances;
DROP TABLE IF EXISTS balances;
DROP FUNCTION IF EXISTS trigger_set_updated_at cascade;
