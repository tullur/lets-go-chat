package memory

import (
	"errors"

	"github.com/tullur/lets-go-chat/internal/domain/chat/token"
)

type MemoryTokenRepository struct {
	Tokens map[string]token.Token
}

func NewMemoryTokenRepositorysitory() *MemoryTokenRepository {
	return &MemoryTokenRepository{
		Tokens: make(map[string]token.Token),
	}
}

func (repo *MemoryTokenRepository) Add(token *token.Token) error {
	if _, ok := repo.Tokens[token.Id()]; ok {
		return errors.New("token already prodvided")
	}

	repo.Tokens[token.Id()] = *token

	return nil
}

func (repo *MemoryTokenRepository) Revoke(id string) error {
	if _, ok := repo.Tokens[id]; !ok {
		return errors.New("token Not Found")
	}

	delete(repo.Tokens, id)

	return nil
}
