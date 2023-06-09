package token

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	id           uuid.UUID
	userId       uuid.UUID
	expiresAfter time.Time
}

func New(userId string) *Token {
	return &Token{id: uuid.New(), userId: uuid.MustParse(userId), expiresAfter: time.Now().Add(time.Hour)}
}

func (t *Token) Id() string {
	return t.id.String()
}

func (t *Token) User() string {
	return t.userId.String()
}

func (t *Token) ExpiresAfter() string {
	return t.expiresAfter.Local().String()
}

func (t *Token) SetID(id string) {
	t.id = uuid.MustParse(id)
}

func (t *Token) SetUser(userId string) {
	t.userId = uuid.MustParse(userId)
}

func (t *Token) SetExpiresAfter(expiresAfter string) {
	expiresTime, err := time.Parse(time.Now().String(), expiresAfter)
	if err != nil {
		return
	}

	t.expiresAfter = expiresTime
}
