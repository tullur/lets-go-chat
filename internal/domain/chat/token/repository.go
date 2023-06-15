package token

type Repository interface {
	Get(id string) (*Token, error)
	Add(token *Token) error
	Revoke(id string) error
}
