package http

import (
	"encoding/json"
	"fmt"
	"github.com/shii-cchi/todo-api/internal/delivery/http/dto"
	"net/http"
)

func (h *Handler) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := new(dto.UserDto)
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := h.todoService.CreateUser(r.Context(), newUser.Name)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Coudn't create user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUsertoUser(user))
}
