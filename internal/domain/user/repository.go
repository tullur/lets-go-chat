package user

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
)

//go:generate mockgen -source=repository.go -destination=./mocks/mock_user_repository.go -package=mocks
type Repository interface {
	List() []User
	Create(user *User) error
	FindByName(userName string) (*User, error)
}
