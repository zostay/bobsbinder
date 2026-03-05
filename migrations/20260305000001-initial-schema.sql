-- +migrate Up
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE parties (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    relationship ENUM('self','spouse','dependent','other') NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE document_categories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INT DEFAULT 0
);

CREATE TABLE documents (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    party_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    status ENUM('draft','complete') DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES document_categories(id)
);

CREATE TABLE document_files (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    document_id BIGINT NOT NULL,
    filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE CASCADE
);

CREATE TABLE checklist_items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    party_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY (party_id, category_id),
    FOREIGN KEY (party_id) REFERENCES parties(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES document_categories(id)
);

CREATE TABLE shares (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    unlock_code_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Seed document categories
INSERT INTO document_categories (slug, name, description, sort_order) VALUES
('instructions', 'Instructions to Survivors', 'Helpful instructions and guidance for your loved ones', 1),
('will', 'Will', 'Last will and testament', 2),
('poa-medical', 'Medical Power of Attorney', 'Medical power of attorney designation', 3),
('poa-financial', 'Financial Power of Attorney', 'Financial power of attorney designation', 4),
('medical-directives', 'Medical Directives', 'Living will, right to death, and other medical directives', 5),
('final-arrangements', 'Final Arrangements', 'Funeral plans, bequests, and final wishes', 6),
('pre-certification', 'Pre-Certification Letters', 'Pre-certification documentation', 7),
('memorial', 'Order of Memorial', 'Memorial service plans and preferences', 8),
('remains', 'Disposition of Remains', 'Instructions for disposition of remains', 9),
('assets', 'Asset Distribution', 'Distribution of assets and property', 10),
('safe-access', 'Safe Combinations & Access', 'Safe combinations, lockbox instructions, and access codes', 11),
('accounts', 'List of Accounts', 'Bank accounts, investment accounts, insurance policies', 12),
('property', 'Properties & Tangible Assets', 'Real estate, vehicles, valuable personal property', 13),
('passwords', 'Password Recovery', 'Password manager recovery codes and digital access', 14),
('other', 'Additional Documents', 'Any other important documents', 15);

-- +migrate Down
DROP TABLE IF EXISTS shares;
DROP TABLE IF EXISTS checklist_items;
DROP TABLE IF EXISTS document_files;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS document_categories;
DROP TABLE IF EXISTS parties;
DROP TABLE IF EXISTS users;
