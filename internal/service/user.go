package service

import (
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/pkg/hasher"
)

type UserService struct {
	Repo user.Repository
}

func (u *UserService) GetList() []user.User {
	return u.Repo.List()
}

func (u *UserService) CreateUser(name, password string) (user.User, error) {
	if name == "" || password == "" {
		return user.User{}, user.ErrEmptyValues
	}

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		return user.User{}, err
	}

	user, err := user.NewUser(name, hashedPassword)
	if err != nil {
		return user, err
	}

	err = u.Repo.Create(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserService) LoginUser(name, password string) (user.User, error) {
	currentUser, err := u.Repo.FindByName(name)
	if err != nil {
		return user.User{}, err
	}

	checker := hasher.CheckPasswordHash(password, currentUser.Password)
	if !checker {
		return user.User{}, user.ErrInvalidPassword
	}

	return currentUser, nil
}
