package authn

type Authn interface {
	Options() AuthnOptions
	Generate(id string, opts ...GenerateOption) (*Account, error)
	String() string
}
