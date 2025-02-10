-- +goose Up
CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    role VARCHAR(50) DEFAULT 'USER',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE auth;
