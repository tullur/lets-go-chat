package user

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
)

type Repository interface {
	List() []User
	Create(user User) error
	FindByName(userName string) (User, error)
}
