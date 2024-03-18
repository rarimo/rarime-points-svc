-- +migrate Up
ALTER TABLE balances
    ADD level integer not null default 0;

ALTER TABLE balances
    ADD level_claim_id text unique;

-- +migrate Down