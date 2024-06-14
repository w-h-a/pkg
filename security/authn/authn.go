package authn

import "errors"

var (
	ErrUserEmailInUse = errors.New("there is already a user with this email")
)

type Authn interface {
	Options() AuthnOptions
	Generate(id string, opts ...GenerateOption) (*Account, error)
	String() string
}
