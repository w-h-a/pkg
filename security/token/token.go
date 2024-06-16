package token

import "errors"

var (
	ErrInvalidToken = errors.New("invalid token provided")
)

type TokenProvider interface {
	Options() TokenOptions
	Generate(opts ...GenerateOption) (*Token, error)
	Inspect(tok string) (*Token, error)
	String() string
}
