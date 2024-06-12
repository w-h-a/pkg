package authn

type Account struct {
	Id       string            `json:"id"`
	Secret   string            `json:"secret"`
	Roles    []string          `json:"roles"`
	Metadata map[string]string `json:"metadata"`
}
