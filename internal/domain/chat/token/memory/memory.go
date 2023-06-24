package memory

import (
	"errors"

	"github.com/tullur/lets-go-chat/internal/domain/chat/token"
)

var (
	ErrTokenNotExists       = errors.New("token does not exist")
	ErrTokenAlreadyProvided = errors.New("token already prodvided")
	ErrTokeNotFound         = errors.New("token Not Found")
)

type MemoryTokenRepository struct {
	Tokens map[string]token.Token
}

func NewMemoryTokenRepository() *MemoryTokenRepository {
	return &MemoryTokenRepository{
		Tokens: make(map[string]token.Token),
	}
}

func (repo *MemoryTokenRepository) Get(id string) (*token.Token, error) {
	if token, ok := repo.Tokens[id]; ok {
		return &token, nil
	}

	return nil, ErrTokenNotExists
}

func (repo *MemoryTokenRepository) Add(token *token.Token) error {
	if _, ok := repo.Tokens[token.Id()]; ok {
		return ErrTokenAlreadyProvided
	}

	repo.Tokens[token.Id()] = *token

	return nil
}

func (repo *MemoryTokenRepository) Revoke(id string) error {
	if _, ok := repo.Tokens[id]; !ok {
		return ErrTokeNotFound
	}

	delete(repo.Tokens, id)

	return nil
}
