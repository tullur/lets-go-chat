package service

import (
	"github.com/tullur/lets-go-chat/internal/domain/chat/token"
	"github.com/tullur/lets-go-chat/internal/domain/chat/token/memory"
	"github.com/tullur/lets-go-chat/internal/domain/user"
)

type ChatConfiguration func(cs *ChatService) error

type ChatService struct {
	tokenRepo token.Repository
}

func WithInMemoryTokenRepository() ChatConfiguration {
	repository := memory.NewMemoryTokenRepository()

	return func(cs *ChatService) error {
		cs.tokenRepo = repository
		return nil
	}
}

func NewTokenService(cfgs ...ChatConfiguration) (*ChatService, error) {
	cs := &ChatService{}

	for _, cfg := range cfgs {
		err := cfg(cs)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}

func (cs *ChatService) GetToken(id string) (*token.Token, error) {
	token, err := cs.tokenRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return token, nil
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
