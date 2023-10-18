package models

import "github.com/google/uuid"

type Todo struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
	UserID uuid.UUID `json:"userId"`
}
