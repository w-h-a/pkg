package authn

import (
	"errors"

	"github.com/w-h-a/pkg/security/token"
)

var (
	ErrUserEmailInUse  = errors.New("there is already a user with this email")
	ErrIncorrectSecret = errors.New("incorrect secret")
)

type Authn interface {
	Options() AuthnOptions
	Generate(id string, opts ...GenerateOption) (*Account, error)
	Token(opts ...TokenOption) (*token.Token, error)
	Inspect(token string) (*Account, error)
	List() ([]*Account, error)
	String() string
}
