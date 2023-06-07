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

func (t *Token) ExpiresAfter() string {
	return t.expiresAfter.Local().String()
}
