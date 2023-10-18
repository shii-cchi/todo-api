package models

import "github.com/google/uuid"

type User struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	ApiKey string    `json:"apiKey"`
}
