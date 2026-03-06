-- +migrate Up
ALTER TABLE contacts ADD COLUMN is_primary BOOLEAN NOT NULL DEFAULT FALSE AFTER notes;

-- +migrate Down
ALTER TABLE contacts DROP COLUMN is_primary;
