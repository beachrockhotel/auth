-- +goose Up
CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    role INTEGER DEFAULT 1,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL DEFAULT NULL
);
-- +goose Down
DROP TABLE auth;