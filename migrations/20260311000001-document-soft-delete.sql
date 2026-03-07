-- +migrate Up
ALTER TABLE documents ADD COLUMN deleted_at TIMESTAMP NULL DEFAULT NULL AFTER updated_at;

-- +migrate Down
ALTER TABLE documents DROP COLUMN deleted_at;
