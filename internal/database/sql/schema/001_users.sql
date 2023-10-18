-- +goose Up

DROP EXTENSION IF EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex'))
);

-- +goose Down

DROP TABLE users;