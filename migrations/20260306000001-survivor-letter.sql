-- +migrate Up

CREATE TABLE contacts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    relationship VARCHAR(100),
    role VARCHAR(100),
    phone VARCHAR(50),
    email VARCHAR(255),
    address TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE locations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    type ENUM('physical','digital') NOT NULL,
    description TEXT,
    address TEXT,
    access_instructions TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE digital_access (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type ENUM('computer','phone','password_manager') NOT NULL,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    instructions TEXT,
    location_id BIGINT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

CREATE TABLE insurance_policies (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    party_id BIGINT NULL,
    provider VARCHAR(255) NOT NULL,
    policy_number VARCHAR(255),
    type VARCHAR(100),
    coverage_amount DECIMAL(15,2),
    beneficiary VARCHAR(255),
    agent_name VARCHAR(255),
    agent_phone VARCHAR(50),
    location_id BIGINT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE SET NULL,
    FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

CREATE TABLE service_accounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    type ENUM('financial_tool','backup_service','tax_preparer') NOT NULL,
    name VARCHAR(255) NOT NULL,
    provider VARCHAR(255),
    account_number VARCHAR(255),
    contact_name VARCHAR(255),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE party_obituary_info (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    party_id BIGINT NOT NULL,
    type ENUM('survivor','predeceased','event') NOT NULL,
    name VARCHAR(255) NOT NULL,
    relationship VARCHAR(100),
    details TEXT,
    event_date DATE NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE
);

CREATE TABLE survivor_letters (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    greeting TEXT,
    intro TEXT,
    closing TEXT,
    signature TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE survivor_letter_sections (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    letter_id BIGINT NOT NULL,
    section_key VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    visible BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (letter_id, section_key),
    FOREIGN KEY (letter_id) REFERENCES survivor_letters(id) ON DELETE CASCADE
);

CREATE TABLE survivor_letter_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    section_id BIGINT NOT NULL,
    source_type VARCHAR(100),
    source_id BIGINT NULL,
    content TEXT NOT NULL,
    item_type ENUM('numbered','bulleted','paragraph') NOT NULL DEFAULT 'numbered',
    provenance ENUM('auto','auto_edited','manual') NOT NULL DEFAULT 'auto',
    suppressed BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (section_id, source_type, source_id),
    FOREIGN KEY (section_id) REFERENCES survivor_letter_sections(id) ON DELETE CASCADE
);

ALTER TABLE documents ADD COLUMN location_id BIGINT NULL;
ALTER TABLE documents ADD CONSTRAINT fk_documents_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL;

-- +migrate Down
ALTER TABLE documents DROP FOREIGN KEY fk_documents_location;
ALTER TABLE documents DROP COLUMN location_id;

DROP TABLE IF EXISTS survivor_letter_items;
DROP TABLE IF EXISTS survivor_letter_sections;
DROP TABLE IF EXISTS survivor_letters;
DROP TABLE IF EXISTS party_obituary_info;
DROP TABLE IF EXISTS service_accounts;
DROP TABLE IF EXISTS insurance_policies;
DROP TABLE IF EXISTS digital_access;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS contacts;
