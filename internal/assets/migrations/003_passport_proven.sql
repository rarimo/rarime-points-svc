-- +migrate Up
ALTER TABLE balances ADD COLUMN is_passport_proven boolean NOT NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE balances DROP COLUMN is_passport_proven;
