-- +goose Up

CREATE TABLE todos (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    status TEXT NOT NULL
);

-- +goose Down

DROP TABLE todos;