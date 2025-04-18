-- +migrate Up

CREATE TABLE IF NOT EXISTS face_event_balances
(
    nullifier  TEXT   PRIMARY KEY NOT NULL REFERENCES balances (nullifier),
    amount     integer NOT NULL,
    created_at integer NOT NULL default EXTRACT('EPOCH' FROM NOW())
);

CREATE INDEX IF NOT EXISTS face_event_balances_nullifier_index ON face_event_balances (nullifier);


INSERT INTO face_event_balances (nullifier, amount, created_at)
SELECT nullifier, 0, EXTRACT('EPOCH' FROM NOW())::integer
FROM balances;

-- +migrate Down

DROP INDEX IF EXISTS face_event_balances_nullifier_index;

DROP TABLE IF EXISTS face_event_balances;

