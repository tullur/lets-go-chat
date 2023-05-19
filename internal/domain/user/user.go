package user

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyValues     = errors.New("empty username or password")
	ErrInvalidPassword = errors.New("invalid password hash")
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"userName"`
	Password string    `json:"password"`
}

func NewUser(name, password string) (User, error) {
	if name == "" || password == "" {
		return User{}, ErrEmptyValues
	}

	return User{
		ID:       uuid.New(),
		Name:     name,
		Password: password,
	}, nil
}
