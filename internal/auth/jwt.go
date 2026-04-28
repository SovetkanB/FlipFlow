package auth

import "github.com/google/uuid"

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
}
