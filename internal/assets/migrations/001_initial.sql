-- +migrate Up
DROP FUNCTION IF EXISTS trigger_set_updated_at cascade;

CREATE FUNCTION trigger_set_updated_at() RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$ BEGIN NEW.updated_at = NOW() at time zone 'utc'; RETURN NEW; END; $$;

CREATE TABLE IF NOT EXISTS balances
(
    id serial primary key,
    did text not null unique,
    amount integer not null default 0,
    updated_at timestamp without time zone not null default NOW()
);

CREATE INDEX IF NOT EXISTS balances_did_index on balances using btree (did);

CREATE TRIGGER set_updated_at
    before update
    on balances
    for each row
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS events
(
    id serial primary key,
    did text not null,
    type text not null,
    is_claimed boolean not null,
    created_at timestamp without time zone not null default NOW()
);

CREATE INDEX IF NOT EXISTS events_did_index on events using btree (did);

-- +migrate Down
DROP INDEX IF EXISTS events_did_index;
DROP TABLE IF EXISTS events;

DROP TRIGGER IF EXISTS set_updated_at on balances;
DROP INDEX IF EXISTS balances_did_index;
DROP TABLE IF EXISTS balances;

DROP FUNCTION IF EXISTS trigger_set_updated_at cascade;
