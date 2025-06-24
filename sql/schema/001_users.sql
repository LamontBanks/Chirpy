-- Goose for database migrations: https://github.com/pressly/goose

-- +goose Up
CREATE TABLE users (
    id uuid                 PRIMARY KEY,
    created_at timestamp    NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
    email TEXT              UNIQUE
                            NOT NULL
);

-- +goose Down
DROP TABLE users;
