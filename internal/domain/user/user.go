package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tullur/lets-go-chat/pkg/hasher"
)

var (
	ErrEmptyValues         = errors.New("empty username or password")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrShortUserName       = errors.New("user name is too short (minimum is 4 characters)")
	ErrShortPasswordLength = errors.New("password is too short (minimum is 8 characters)")
)

const (
	userNameLength = 4
	passwordLength = 8
)

type User struct {
	id       uuid.UUID
	name     string
	password string
}

func New(name, password string) (*User, error) {
	user := &User{id: uuid.New(), name: name, password: password}

	if err := user.validate(); err != nil {
		return nil, err
	}

	if err := user.hashPassword(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) Id() string {
	return u.id.String()
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Password() string {
	return u.password
}

func (u *User) SetID(id uuid.UUID) {
	u.id = id
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) VerifyPassword(password string) (bool, error) {
	verified := hasher.CheckPasswordHash(password, u.password)
	if !verified {
		return false, ErrInvalidPassword
	}

	return true, nil
}

func (u *User) hashPassword() error {
	hashedPassword, err := hasher.HashPassword(u.password)
	if err != nil {
		return err
	}

	u.password = hashedPassword

	return nil
}

func (u *User) validate() error {
	if u.name == "" || u.password == "" {
		return ErrEmptyValues
	}

	if len(u.name) < userNameLength {
		return ErrShortUserName
	}

	if len(u.password) < passwordLength {
		return ErrShortPasswordLength
	}

	return nil
}
