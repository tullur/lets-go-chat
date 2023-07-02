package memory

import (
	"sync"

	"github.com/tullur/lets-go-chat/internal/domain/user"
)

type InMemoryRepository struct {
	users map[string]user.User
	sync.Mutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[string]user.User),
	}
}

func (repo *InMemoryRepository) List() []user.User {
	var users []user.User

	for _, user := range repo.users {
		users = append(users, user)
	}

	return users
}

func (repo *InMemoryRepository) Create(u *user.User) error {
	repo.Lock()
	defer repo.Unlock()

	if _, ok := repo.users[u.Id()]; ok {
		return user.ErrUserAlreadyExists
	}

	repo.users[u.Id()] = *u

	return nil
}

func (repo *InMemoryRepository) FindByName(name string) (*user.User, error) {
	for _, v := range repo.users {
		if v.Name() == name {
			return &v, nil
		}
	}

	return nil, user.ErrUserNotFound
}
