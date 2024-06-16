package jsonwebtoken

import "github.com/golang-jwt/jwt"

type AuthClaims struct {
	Roles    []string          `json:"roles"`
	Metadata map[string]string `json:"metadata"`

	jwt.StandardClaims
}
