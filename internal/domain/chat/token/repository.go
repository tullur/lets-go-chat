package token

type Repository interface {
	Add(token *Token) error
	Revoke(id string) error
}
