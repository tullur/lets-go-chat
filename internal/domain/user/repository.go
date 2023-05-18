package user

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
)

type Repository interface {
	List() []User
	Create(user User) User
	FindByName(userName string) (User, error)
}

type InMemoryRepo struct {
	Users []User
}

func New() *InMemoryRepo {
	return &InMemoryRepo{
		Users: []User{},
	}
}

func (repo *InMemoryRepo) List() []User {
	return repo.Users
}

func (repo *InMemoryRepo) Create(user User) User {
	repo.Users = append(repo.Users, user)

	return user
}

func (repo *InMemoryRepo) FindByName(userName string) (User, error) {
	for _, v := range repo.Users {
		if v.Name == userName {
			return v, nil
		}
	}

	return User{}, ErrUserNotFound
}
