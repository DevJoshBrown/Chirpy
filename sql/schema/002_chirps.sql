-- +goose Up
CREATE TABLE chirps (
    id uuid primary key,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id uuid NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );

-- +goose Down
DROP TABLE chirps;
