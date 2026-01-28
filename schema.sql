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

-- Indexes (Optional for performance)
CREATE INDEX idx_status ON pdf_files(status);
CREATE INDEX idx_created_at ON pdf_files(created_at);
