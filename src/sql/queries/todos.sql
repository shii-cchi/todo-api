-- name: CreateTodo :one
INSERT INTO todos (id, title, status, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTodosList :many
SELECT id, title, status, user_id FROM todos;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;

-- name: GetTodo :one
SELECT id, title, status, user_id FROM todos WHERE id = $1;

-- name: UpdateTodo :one
UPDATE todos
SET title = $2, status = $3
WHERE id = $1
RETURNING *;