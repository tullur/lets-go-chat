package token

//go:generate mockgen -source=repository.go -destination=./mocks/mock_token_repository.go -package=mocks
type Repository interface {
	Get(id string) (*Token, error)
	Add(token *Token) error
	Revoke(id string) error
}
