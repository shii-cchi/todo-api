-- name: CreateTodo :one
INSERT INTO todos (title, status, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTodosList :many
SELECT id, title, status, user_id FROM todos
WHERE user_id = $1;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1 AND user_id = $2;

-- name: GetTodo :one
SELECT id, title, status, user_id FROM todos
WHERE id = $1 AND user_id = $2;

-- name: UpdateTodo :one
UPDATE todos
SET title = $3, status = $4
WHERE id = $1 AND user_id = $2
RETURNING *;