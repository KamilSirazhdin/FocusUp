-- migrations/YYYYMMDDHHMMSS_add_password_reset_table.sql
-- +goose Up
CREATE TABLE IF NOT EXISTS password_resets (
                                               id SERIAL PRIMARY KEY,
                                               email VARCHAR(100) NOT NULL,
    token VARCHAR(255) NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
    );

CREATE INDEX IF NOT EXISTS idx_password_resets_email ON password_resets(email);
CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token);

-- +goose Down
DROP TABLE IF EXISTS password_resets;