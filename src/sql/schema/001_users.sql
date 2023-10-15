-- +goose Up

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex'))
);

-- +goose Down

DROP TABLE users;