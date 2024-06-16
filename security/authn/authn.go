package authn

import (
	"errors"

	"github.com/w-h-a/pkg/security/token"
)

var (
	ErrUserEmailInUse = errors.New("there is already a user with this email")
)

type Authn interface {
	Options() AuthnOptions
	Generate(id string, opts ...GenerateOption) (*Account, error)
	Token(opts ...TokenOption) (*token.Token, error)
	String() string
}
