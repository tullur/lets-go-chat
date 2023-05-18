package user

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"userName"`
	Password string    `json:"password"`
}
