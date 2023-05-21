package user

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyValues         = errors.New("empty username or password")
	ErrInvalidPassword     = errors.New("invalid password hash")
	ErrShortUserName       = errors.New("user name is too short (minimum is 4 characters)")
	ErrShortPasswordLength = errors.New("password is too short (minimum is 8 characters)")
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"userName"`
	Password string    `json:"password"`
}

func New(name, password string) (User, error) {
	if name == "" || password == "" {
		return User{}, ErrEmptyValues
	}

	return User{
		ID:       uuid.New(),
		Name:     name,
		Password: password,
	}, nil
}
