package token

import "time"

type Token struct {
	AccessToken string            `json:"access_token"`
	Created     time.Time         `json:"created"`
	Expiry      time.Time         `json:"expiry"`
	Id          string            `json:"id"`
	Roles       []string          `json:"roles"`
	Metadata    map[string]string `json:"metadata"`
}
