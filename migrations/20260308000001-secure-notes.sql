-- +migrate Up
ALTER TABLE contacts ADD COLUMN secure_notes TEXT AFTER is_primary;
ALTER TABLE locations ADD COLUMN secure_notes TEXT AFTER access_instructions;
ALTER TABLE digital_access ADD COLUMN secure_notes TEXT AFTER instructions;
ALTER TABLE service_accounts ADD COLUMN secure_notes TEXT AFTER notes;
ALTER TABLE insurance_policies ADD COLUMN secure_notes TEXT AFTER notes;
ALTER TABLE documents ADD COLUMN secure_notes TEXT AFTER status;

-- +migrate Down
ALTER TABLE contacts DROP COLUMN secure_notes;
ALTER TABLE locations DROP COLUMN secure_notes;
ALTER TABLE digital_access DROP COLUMN secure_notes;
ALTER TABLE service_accounts DROP COLUMN secure_notes;
ALTER TABLE insurance_policies DROP COLUMN secure_notes;
ALTER TABLE documents DROP COLUMN secure_notes;
