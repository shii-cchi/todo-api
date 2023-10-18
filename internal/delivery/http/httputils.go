package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shii-cchi/todo-api/internal/database"
	"github.com/shii-cchi/todo-api/internal/models"
	"log"
	"net/http"
)

func (h *Handler) CheckAuth(w http.ResponseWriter, r *http.Request) uuid.UUID {
	apiKey := r.Header.Get("Authorization")

	if apiKey == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return uuid.Nil
	}

	userId, err := h.todoService.CheckAuth(r.Context(), apiKey)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't check: %v", err))
		return uuid.Nil
	}

	if userId == uuid.Nil {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return uuid.Nil
	}

	return userId
}

func (h *Handler) GetId(w http.ResponseWriter, r *http.Request) uuid.UUID {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't parse id: %v", err))
		return uuid.Nil
	}

	return id
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON responce: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func databaseTodotoTodo(dbTodo database.Todo) models.Todo {
	return models.Todo{
		ID:     dbTodo.ID,
		Title:  dbTodo.Title,
		Status: dbTodo.Status,
		UserID: dbTodo.UserID,
	}
}

func databaseUsertoUser(dbTodo database.User) models.User {
	return models.User{
		ID:     dbTodo.ID,
		Name:   dbTodo.Name,
		ApiKey: dbTodo.ApiKey,
	}
}

func databaseTodoListtoTodoList(dbTodos []database.Todo) []models.Todo {
	todos := make([]models.Todo, len(dbTodos))
	for i, dbTodo := range dbTodos {
		todos[i] = databaseTodotoTodo(dbTodo)
	}
	return todos
}
