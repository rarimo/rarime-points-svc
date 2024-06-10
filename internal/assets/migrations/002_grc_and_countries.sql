-- +migrate Up
ALTER TABLE referrals
    RENAME COLUMN is_consumed TO usage_left;
ALTER TABLE referrals
    ALTER COLUMN usage_left DROP DEFAULT;

ALTER TABLE referrals
    ALTER usage_left TYPE INTEGER
        USING
        CASE
            WHEN usage_left = TRUE THEN 0
            ELSE 1
            END;

ALTER TABLE referrals
    ALTER COLUMN usage_left SET DEFAULT 1;
ALTER TABLE balances
    DROP CONSTRAINT balances_referred_by_key;

CREATE TABLE IF NOT EXISTS countries
(
    code               text PRIMARY KEY,
    reserve_limit      bigint  NOT NULL,
    reserved           bigint  NOT NULL default 0,
    withdrawn          bigint  NOT NULL default 0,
    reserve_allowed    boolean NOT NULL,
    withdrawal_allowed boolean NOT NULL
);

ALTER TABLE balances
    ADD COLUMN country text REFERENCES countries (code);

-- +migrate Down
ALTER TABLE balances
    DROP COLUMN country;

DROP TABLE IF EXISTS countries;

ALTER TABLE referrals
    ALTER COLUMN usage_left DROP DEFAULT;

ALTER TABLE referrals
    ALTER usage_left TYPE BOOLEAN
        USING
        CASE
            WHEN usage_left > 0 THEN FALSE
            ELSE TRUE
            END;

ALTER TABLE referrals
    RENAME COLUMN usage_left TO is_consumed;
ALTER TABLE referrals
    ALTER COLUMN is_consumed SET DEFAULT FALSE;
