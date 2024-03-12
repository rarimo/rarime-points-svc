-- +migrate Up
CREATE OR REPLACE FUNCTION trigger_set_updated_at() RETURNS trigger
    LANGUAGE plpgsql
AS $$ BEGIN NEW.updated_at = EXTRACT('EPOCH' FROM NOW()); RETURN NEW; END; $$;

CREATE TABLE IF NOT EXISTS balances
(
    did              text PRIMARY KEY,
    amount           bigint NOT NULL default 0,
    created_at       integer NOT NULL default EXTRACT('EPOCH' FROM NOW()),
    updated_at       integer NOT NULL default EXTRACT('EPOCH' FROM NOW()),
    referred_by      text UNIQUE,
    passport_hash    text UNIQUE,
    passport_expires timestamp without time zone
    is_withdrawal_allowed boolean NOT NULL default true
);

CREATE INDEX IF NOT EXISTS balances_page_index ON balances (amount, updated_at) WHERE referred_by IS NOT NULL;

CREATE TRIGGER set_updated_at
    BEFORE UPDATE
    ON balances
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS referrals
(
    id          text PRIMARY KEY,
    user_did    text    NOT NULL REFERENCES balances (did),
    is_consumed boolean NOT NULL default false
);

ALTER TABLE balances ADD CONSTRAINT referred_by_fk FOREIGN KEY (referred_by) REFERENCES referrals (id);
CREATE INDEX IF NOT EXISTS referrals_user_did_index ON referrals (user_did);

CREATE TYPE event_status AS ENUM ('open', 'fulfilled', 'claimed');

CREATE TABLE IF NOT EXISTS events
(
    id            uuid PRIMARY KEY NOT NULL default gen_random_uuid(),
    user_did      text             NOT NULL REFERENCES balances (did),
    type          text             NOT NULL,
    status        event_status     NOT NULL,
    created_at    integer          NOT NULL default EXTRACT('EPOCH' FROM NOW()),
    updated_at    integer          NOT NULL default EXTRACT('EPOCH' FROM NOW()),
    meta          jsonb,
    points_amount integer,
    external_id   text,
    CONSTRAINT unique_external_id UNIQUE (user_did, type, external_id)
);

CREATE INDEX IF NOT EXISTS events_page_index ON events (user_did, updated_at);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE
    ON events
    FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE IF NOT EXISTS withdrawals
(
    id         uuid PRIMARY KEY default gen_random_uuid(),
    user_did   text    NOT NULL REFERENCES balances (did),
    amount     integer NOT NULL,
    address    text    NOT NULL,
    created_at integer NOT NULL default EXTRACT('EPOCH' FROM NOW())
);

CREATE INDEX IF NOT EXISTS withdrawals_page_index ON withdrawals (user_did, created_at);

-- +migrate Down
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS events;

ALTER TABLE balances DROP CONSTRAINT referred_by_fk;
DROP TABLE IF EXISTS referrals;
DROP TABLE IF EXISTS balances;

DROP TYPE IF EXISTS event_status;
DROP FUNCTION IF EXISTS trigger_set_updated_at();
