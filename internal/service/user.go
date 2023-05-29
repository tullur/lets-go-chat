package service

import "github.com/tullur/lets-go-chat/internal/domain/user"

type UserService struct {
	repository user.Repository
}

func (u *UserService) GetList() []user.User {
	return u.repository.List()
}

func (u *UserService) CreateUser(name, password string) (*user.User, error) {
	user, err := user.New(name, password)
	if err != nil {
		return nil, err
	}

	err = u.repository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) LoginUser(name, password string) (user.User, error) {
	currentUser, err := u.repository.FindByName(name)
	if err != nil {
		return user.User{}, err
	}

	verified, err := currentUser.VerifyPassword(password)
	if !verified {
		return user.User{}, err
	}
	return currentUser, nil
}

func NewUserService(repository user.Repository) *UserService {
	return &UserService{repository: repository}
}
