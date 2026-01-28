-- Database Schema for PDF Management System

CREATE TABLE IF NOT EXISTS pdf_files (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255),
    filepath VARCHAR(500) NOT NULL,
    size BIGINT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('CREATED', 'UPLOADED', 'DELETED')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    role VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    address VARCHAR(255),
    created_by BIGINT,
    created_date TIMESTAMP,
    email VARCHAR(30) UNIQUE NOT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    modified_by BIGINT,
    modified_date TIMESTAMP,
    name VARCHAR(50),
    password VARCHAR(255) NOT NULL,
    phone_number VARCHAR(13),
    post_code CHAR(5),
    role_id BIGINT REFERENCES roles(id)
);

-- Seed Initial Roles
INSERT INTO roles (role) VALUES ('Project Manager'), ('Financial'), ('HRD') ON CONFLICT DO NOTHING;

-- Indexes (Optional for performance)
CREATE INDEX idx_status ON pdf_files(status);
CREATE INDEX idx_created_at ON pdf_files(created_at);
