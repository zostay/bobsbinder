-- +migrate Up
ALTER TABLE checklist_items
  ADD COLUMN status ENUM('pending','complete','not_applicable') NOT NULL DEFAULT 'pending' AFTER category_id;
UPDATE checklist_items SET status = 'complete' WHERE completed = TRUE;
ALTER TABLE checklist_items DROP COLUMN completed;

-- +migrate Down
ALTER TABLE checklist_items ADD COLUMN completed BOOLEAN DEFAULT FALSE AFTER category_id;
UPDATE checklist_items SET completed = TRUE WHERE status = 'complete';
ALTER TABLE checklist_items DROP COLUMN status;
