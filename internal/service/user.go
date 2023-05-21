package service

import (
	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/pkg/hasher"
)

type UserService struct {
	repo user.Repository
}

func NewUser(repository user.Repository) *UserService {
	return &UserService{
		repo: repository,
	}
}

func (u *UserService) GetList() []user.User {
	return u.repo.List()
}

func (u *UserService) Create(name, password string) (user.User, error) {
	if name == "" || password == "" {
		return user.User{}, user.ErrEmptyValues
	}

	if len(name) < 4 {
		return user.User{}, user.ErrShortUserName
	}

	if len(password) < 8 {
		return user.User{}, user.ErrShortPasswordLength
	}

	hashedPassword, err := hasher.HashPassword(password)
	if err != nil {
		return user.User{}, err
	}

	user, err := user.New(name, hashedPassword)
	if err != nil {
		return user, err
	}

	u.repo.Create(user)

	return user, nil
}

func (u *UserService) Login(name, password string) (user.User, error) {
	currentUser, err := u.repo.FindByName(name)
	if err != nil {
		return user.User{}, err
	}

	checker := hasher.CheckPasswordHash(password, currentUser.Password)
	if !checker {
		return user.User{}, user.ErrInvalidPassword
	}

	return currentUser, nil
}
