package authz

import "errors"

var (
	ErrForbidden = errors.New("resource forbidden")
)

type Authz interface {
	Options() AuthzOptions
	Grant(role string, res *Resource) error
	Revoke(role string, res *Resource) error
	Verify(roles []string, res *Resource) error
	List() ([]*Rule, error)
	String() string
}
