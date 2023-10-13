package main

import (
	"github.com/google/uuid"
	"github.com/shii-cchi/todo-api/internal/database"
)

type Todo struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
}

func databaseTodotoTodo(dbTodo database.Todo) Todo {
	return Todo{
		ID:     dbTodo.ID,
		Title:  dbTodo.Title,
		Status: dbTodo.Status,
	}
}

func databaseTodoListtoTodoList(dbTodos []database.Todo) []Todo {
	todos := make([]Todo, len(dbTodos))
	for i, dbTodo := range dbTodos {
		todos[i] = databaseTodotoTodo(dbTodo)
	}
	return todos
}
