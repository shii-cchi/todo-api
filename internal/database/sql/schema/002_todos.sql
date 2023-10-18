-- +goose Up

CREATE TABLE todos (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    title TEXT NOT NULL,
    status TEXT NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down

DROP TABLE todos;