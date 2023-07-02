package user

import "errors"

var (
	ErrUserNotFound      = errors.New("User not found")
	ErrUserAlreadyExists = errors.New("User already exists")
)

//go:generate mockgen -source=repository.go -destination=./mocks/mock_user_repository.go -package=mocks
type Repository interface {
	List() []User
	Create(user *User) error
	FindByName(name string) (*User, error)
}
