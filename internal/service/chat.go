package service

import (
	"github.com/tullur/lets-go-chat/internal/domain/chat/token"
	"github.com/tullur/lets-go-chat/internal/domain/user"
)

type ChatConfiguration func(cs *ChatService) error

type ChatService struct {
	tokenRepo token.Repository
}

func NewTokenService(tokeRepo token.Repository) *ChatService {
	return &ChatService{tokenRepo: tokeRepo}
}

func (cs *ChatService) GenerateAccessToken(user *user.User) (*token.Token, error) {
	token := token.New(user.Id())

	err := cs.tokenRepo.Add(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (cs *ChatService) RevokeToken(id string) error {
	err := cs.tokenRepo.Revoke(id)
	if err != nil {
		return err
	}

	return nil
}
