package http

import (
	"encoding/json"
	"fmt"
	"github.com/shii-cchi/todo-api/internal/database"
	"github.com/shii-cchi/todo-api/internal/delivery/http/dto"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "my todo list")
}

func (h *Handler) handlerFetchTodos(w http.ResponseWriter, r *http.Request) {
	userId := h.CheckAuth(w, r)

	todoList, err := h.todoService.GetTodosList(r.Context(), userId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't get todos: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodoListtoTodoList(todoList))
}

func (h *Handler) handlerFetchTodo(w http.ResponseWriter, r *http.Request) {
	id := h.GetId(w, r)

	userId := h.CheckAuth(w, r)

	todo, err := h.todoService.GetTodo(r.Context(), database.GetTodoParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't get todo: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodotoTodo(todo))
}

func (h *Handler) handlerCreateTodo(w http.ResponseWriter, r *http.Request) {
	newTodo := new(dto.TodoDto)
	err := json.NewDecoder(r.Body).Decode(&newTodo)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	userId := h.CheckAuth(w, r)

	todo, err := h.todoService.CreateTodo(r.Context(), database.CreateTodoParams{
		Title:  newTodo.Title,
		Status: newTodo.Status,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create todo: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseTodotoTodo(todo))
}

func (h *Handler) handlerUpdateTodo(w http.ResponseWriter, r *http.Request) {
	id := h.GetId(w, r)

	userId := h.CheckAuth(w, r)

	todoOld, err := h.todoService.GetTodo(r.Context(), database.GetTodoParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find todos: %v", err))
		return
	}

	newTodo := new(dto.TodoDto)
	err = json.NewDecoder(r.Body).Decode(&newTodo)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	todoNew, err := h.todoService.UpdateTodo(r.Context(), database.UpdateTodoParams{
		ID:     id,
		Title:  newTodo.Title,
		Status: newTodo.Status,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't update user: %v", err))
		return
	}

	if todoNew.Title == todoOld.Title && todoNew.Status == todoOld.Status {
		respondWithError(w, http.StatusNoContent, fmt.Sprintf("No content: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodotoTodo(todoNew))
}

func (h *Handler) handlerDeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := h.GetId(w, r)

	userId := h.CheckAuth(w, r)

	todo, err := h.todoService.GetTodo(r.Context(), database.GetTodoParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Coudn't find todos: %v", err))
		return
	}

	err = h.todoService.DeleteTodo(r.Context(), database.DeleteTodoParams{
		ID:     id,
		UserID: userId,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't delete: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databaseTodotoTodo(todo))
}
