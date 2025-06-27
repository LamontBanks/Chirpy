-- Goose for database migrations: https://github.com/pressly/goose

-- +goose Up
CREATE TABLE chirps (
    id uuid PRIMARY KEY,
    created_at timestamp    NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp    NOT NULL
                            DEFAULT CURRENT_TIMESTAMP,
    body    TEXT            NOT NULL,
    user_id uuid            NOT NULL 
                            REFERENCES users
                             -- DELETES this row if the user_id in `users` is deleted
                            ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;