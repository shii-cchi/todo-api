package http

import (
	"github.com/shii-cchi/todo-api/internal/database"
	"github.com/shii-cchi/todo-api/internal/service"
)

type Handler struct {
	todoService *service.TodoService
}

func New(queries *database.Queries) *Handler {
	return &Handler{
		todoService: service.NewTodoService(queries),
	}
}
