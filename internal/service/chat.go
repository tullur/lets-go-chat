package service

import (
	"github.com/tullur/lets-go-chat/internal/domain/chat"
	"github.com/tullur/lets-go-chat/internal/domain/user"
)

func GenerateAccessToken(user *user.User) (string, string) {
	token := chat.NewToken(user.Id())

	return token.Id(), token.ExpiresAfter()
}
