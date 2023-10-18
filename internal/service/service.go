package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/shii-cchi/todo-api/internal/database"
)

type TodoService struct {
	queries *database.Queries
}

func NewTodoService(q *database.Queries) *TodoService {
	return &TodoService{
		queries: q,
	}
}

func (ts *TodoService) CheckAuth(ctx context.Context, apiKey string) (uuid.UUID, error) {
	userId, err := ts.queries.CheckApiKey(ctx, apiKey)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func (ts *TodoService) GetTodosList(ctx context.Context, userId uuid.UUID) ([]database.Todo, error) {
	todoList, err := ts.queries.GetTodosList(ctx, userId)
	if err != nil {
		return nil, err
	}
	return todoList, nil
}

func (ts *TodoService) GetTodo(ctx context.Context, params database.GetTodoParams) (database.Todo, error) {
	todo, err := ts.queries.GetTodo(ctx, params)
	if err != nil {
		return database.Todo{}, err
	}
	return todo, nil
}

func (ts *TodoService) CreateTodo(ctx context.Context, params database.CreateTodoParams) (database.Todo, error) {
	todo, err := ts.queries.CreateTodo(ctx, params)
	if err != nil {
		return database.Todo{}, err
	}
	return todo, nil
}

func (ts *TodoService) UpdateTodo(ctx context.Context, params database.UpdateTodoParams) (database.Todo, error) {
	todo, err := ts.queries.UpdateTodo(ctx, params)
	if err != nil {
		return database.Todo{}, err
	}
	return todo, nil
}

func (ts *TodoService) DeleteTodo(ctx context.Context, params database.DeleteTodoParams) error {
	err := ts.queries.DeleteTodo(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TodoService) CreateUser(ctx context.Context, name string) (database.User, error) {
	user, err := ts.queries.CreateUser(ctx, name)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}
