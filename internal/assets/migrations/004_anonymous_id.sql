-- +migrate Up
ALTER TABLE balances ADD COLUMN anonymous_id text UNIQUE;

-- +migrate Down
ALTER TABLE balances DROP COLUMN anonymous_id;
