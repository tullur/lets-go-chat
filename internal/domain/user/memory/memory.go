package memory

import "github.com/tullur/lets-go-chat/internal/domain/user"

type InMemoryRepo struct {
	Users []user.User
}

func NewInMemoryRepository() *InMemoryRepo {
	return &InMemoryRepo{
		Users: []user.User{},
	}
}

func (repo *InMemoryRepo) List() []user.User {
	return repo.Users
}

func (repo *InMemoryRepo) Create(user user.User) {
	repo.Users = append(repo.Users, user)
}

func (repo *InMemoryRepo) FindByName(userName string) (user.User, error) {
	for _, v := range repo.Users {
		if v.Name == userName {
			return v, nil
		}
	}

	return user.User{}, user.ErrUserNotFound
}
