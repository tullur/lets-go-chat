package chat

import "github.com/tullur/lets-go-chat/internal/domain/user"

type Chat struct {
	ActiveUsers  []user.User
	AccessTokens []string
}
