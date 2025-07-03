-- Goose for database migrations: https://github.com/pressly/goose

-- +goose Up
CREATE TABLE refresh_tokens (
    token       TEXT        PRIMARY KEY,
    created_at  timestamp   NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp   NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
    user_id     uuid        NOT NULL
                            REFERENCES users
                            -- DELETE this row if the user_id is deleted in `users`
                            ON DELETE CASCADE,
    expires_at  timestamp   NOT NULL,
    revoked_at  timestamp   -- NULL if token has not been revoked
);

-- +goose Down
DROP TABLE refresh_tokens;