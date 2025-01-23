-- +goose Up
CREATE TABLE auth (
                      id SERIAL PRIMARY KEY,                     -- Уникальный идентификатор пользователя
                      name VARCHAR(255) NOT NULL,                -- Имя пользователя
                      email VARCHAR(255) UNIQUE NOT NULL,        -- Уникальный адрес электронной почты
                      password_hash TEXT NOT NULL,               -- Хэш пароля пользователя
                      role VARCHAR(50) DEFAULT 'USER',           -- Роль пользователя (USER, ADMIN и т.д.)
                      created_at TIMESTAMP NOT NULL DEFAULT NOW(), -- Дата и время создания записи
                      updated_at TIMESTAMP DEFAULT NULL          -- Дата и время последнего обновления записи
);

-- +goose Down
DROP TABLE auth;
