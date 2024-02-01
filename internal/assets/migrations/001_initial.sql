-- +migrate Up
CREATE OR REPLACE FUNCTION trigger_set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
AS $$ BEGIN NEW.updated_at = EXTRACT('EPOCH' FROM NOW()); RETURN NEW; END; $$;

CREATE TABLE IF NOT EXISTS balances
(
    did        text PRIMARY KEY,
    amount     integer not null default 0,
    created_at integer not null default EXTRACT('EPOCH' FROM NOW()),
    updated_at integer not null default EXTRACT('EPOCH' FROM NOW())
);

CREATE INDEX IF NOT EXISTS balances_amount_index ON balances using btree (amount);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE
    ON balances
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TYPE event_status AS ENUM ('open', 'fulfilled', 'claimed');

CREATE TABLE IF NOT EXISTS events
(
    id            uuid PRIMARY KEY not null default gen_random_uuid(),
    user_did      text             not null REFERENCES balances (did),
    type          text             not null,
    status        event_status     not null,
    created_at    integer          not null default EXTRACT('EPOCH' FROM NOW()),
    updated_at    integer          not null default EXTRACT('EPOCH' FROM NOW()),
    meta          jsonb,
    points_amount integer
);

CREATE INDEX IF NOT EXISTS events_user_did_index ON events using btree (user_did);
CREATE INDEX IF NOT EXISTS events_type_index ON events using btree (type);
CREATE INDEX IF NOT EXISTS events_updated_at_index ON events using btree (updated_at);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE
    ON events
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

-- +migrate Down
DROP TRIGGER IF EXISTS set_updated_at ON events;
DROP INDEX IF EXISTS events_type_index;
DROP INDEX IF EXISTS events_user_did_index;
DROP INDEX IF EXISTS events_updated_at_index;
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_status;

DROP TRIGGER IF EXISTS set_updated_at ON balances;
DROP INDEX IF EXISTS balances_amount_index;
DROP TABLE IF EXISTS balances;

DROP FUNCTION IF EXISTS trigger_set_updated_at();
