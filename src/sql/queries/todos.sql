-- name: CreateTodo :one

INSERT INTO todos (id, title, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTodosList :many

SELECT id, title, status FROM todos;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;

-- name: GetTodo :one
SELECT id, title, status FROM todos WHERE id = $1;

-- name: UpdateTodo :one
UPDATE todos
SET title = $2, status = $3
WHERE id = $1
RETURNING *;