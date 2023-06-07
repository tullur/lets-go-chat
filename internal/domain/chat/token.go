package chat

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	id        uuid.UUID
	userId    uuid.UUID
	expiresAt time.Time
}

func NewToken(userId string) *Token {
	return &Token{id: uuid.New(), userId: uuid.MustParse(userId), expiresAt: time.Now().Add(time.Hour)}
}

func (t *Token) Id() string {
	return t.id.String()
}
