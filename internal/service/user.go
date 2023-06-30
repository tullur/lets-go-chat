package service

import (
	"context"

	"github.com/tullur/lets-go-chat/internal/domain/user"
	"github.com/tullur/lets-go-chat/internal/domain/user/memory"
	"github.com/tullur/lets-go-chat/internal/domain/user/mongo"
)

type UserConfiguration func(u *UserService) error

type UserService struct {
	repository user.Repository
}

func NewUserService(cfgs ...UserConfiguration) (*UserService, error) {
	us := &UserService{}

	for _, cfg := range cfgs {
		err := cfg(us)
		if err != nil {
			return nil, err
		}
	}

	return us, nil
}

func WithRepository(ur user.Repository) UserConfiguration {
	return func(u *UserService) error {
		u.repository = ur
		return nil
	}
}

func WithInMemoryRepository() UserConfiguration {
	repository := memory.NewInMemoryRepository()
	return WithRepository(repository)
}

func WithMongoRepository(connection string) UserConfiguration {
	return func(u *UserService) error {
		repository, err := mongo.New(context.Background(), connection)
		if err != nil {
			return err
		}

		u.repository = repository
		return nil
	}
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

func (u *UserService) LoginUser(name, password string) (*user.User, error) {
	currentUser, err := u.repository.FindByName(name)
	if err != nil {
		return nil, err
	}

	verified, err := currentUser.VerifyPassword(password)
	if !verified {
		return nil, err
	}

	return currentUser, nil
}
