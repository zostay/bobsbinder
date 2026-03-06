-- +migrate Up
DELETE sls FROM survivor_letter_sections sls
WHERE sls.section_key IN (
    'primary_documents', 'additional_documents', 'physical_documents',
    'insurance_policies', 'obituary_info',
    'digital_locations', 'other_locations',
    'digital_access', 'password_managers',
    'financial_tools', 'backup_services', 'tax_preparer'
);

-- +migrate Down
DELETE sls FROM survivor_letter_sections sls
WHERE sls.section_key IN ('documents', 'locations', 'digital_info');
