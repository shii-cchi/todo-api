// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const checkApiKey = `-- name: CheckApiKey :one
SELECT id FROM users
WHERE api_key = $1
`

func (q *Queries) CheckApiKey(ctx context.Context, apiKey string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, checkApiKey, apiKey)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, name)
VALUES ($1, $2)
RETURNING id, name, api_key
`

type CreateUserParams struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Name)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.ApiKey)
	return i, err
}