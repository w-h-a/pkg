package authz

import (
	"encoding/json"
	"strings"
)

const (
	RuleJoinKey = ":"
)

type Resource struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint"`
}

type Rule struct {
	Role     string    `json:"rule"`
	Resource *Resource `json:"resource"`
}

func (r *Rule) Key() string {
	components := []string{r.Resource.Namespace, r.Resource.Name, r.Resource.Endpoint, r.Role}
	return strings.Join(components, RuleJoinKey)
}

func (r *Rule) Bytes() []byte {
	bytes, _ := json.Marshal(r)
	return bytes
}
