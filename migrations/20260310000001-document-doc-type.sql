-- +migrate Up
ALTER TABLE documents ADD COLUMN doc_type ENUM('reference','typed') NOT NULL DEFAULT 'reference' AFTER status;

-- +migrate Down
ALTER TABLE documents DROP COLUMN doc_type;
