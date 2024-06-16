package jsonwebtoken

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/w-h-a/pkg/security/token"
)

type jsonWebTokenProvider struct {
	options    token.TokenOptions
	publicKey  string
	privateKey string
}

func (t *jsonWebTokenProvider) Options() token.TokenOptions {
	return t.options
}

func (t *jsonWebTokenProvider) Generate(opts ...token.GenerateOption) (*token.Token, error) {
	// decode the private
	priv, err := base64.StdEncoding.DecodeString(t.privateKey)
	if err != nil {
		return nil, err
	}

	// parse the private
	key, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	// parse options
	options := token.NewGenerateOptions(opts...)

	// get expiry
	expiry := time.Now().Add(options.Expiry)

	// generate the jwt
	jsonWebToken := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		AuthClaims{
			Roles:    options.Roles,
			Metadata: options.Metadata,
			StandardClaims: jwt.StandardClaims{
				Subject:   options.Id,
				ExpiresAt: expiry.Unix(),
			},
		},
	)

	// sign with private key
	signed, err := jsonWebToken.SignedString(key)
	if err != nil {
		return nil, err
	}

	// return token
	tk := &token.Token{
		AccessToken: signed,
		Created:     time.Now(),
		Expiry:      expiry,
		Id:          options.Id,
		Roles:       options.Roles,
		Metadata:    options.Metadata,
	}

	return tk, nil
}

func (t *jsonWebTokenProvider) Inspect(tok string) (*token.Token, error) {
	// decode the public key
	pub, err := base64.StdEncoding.DecodeString(t.publicKey)
	if err != nil {
		return nil, err
	}

	// parse public
	result, err := jwt.ParseWithClaims(
		tok,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwt.ParseRSAPublicKeyFromPEM(pub)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// validate
	if !result.Valid {
		return nil, token.ErrInvalidToken
	}

	// get the claims
	claims, ok := result.Claims.(*AuthClaims)
	if !ok {
		return nil, token.ErrInvalidToken
	}

	// return token as account/claims info
	tk := &token.Token{
		Id:       claims.Subject,
		Roles:    claims.Roles,
		Metadata: claims.Metadata,
	}

	return tk, nil
}

func (t *jsonWebTokenProvider) String() string {
	return "jwt"
}

func NewTokenProvider(opts ...token.TokenOption) token.TokenProvider {
	options := token.NewTokenOptions(opts...)

	j := &jsonWebTokenProvider{
		options: options,
	}

	publicKey, ok := GetPublicKeyFromContext(options.Context)
	if ok {
		j.publicKey = publicKey
	}

	privateKey, ok := GetPrivateKeyFromContext(options.Context)
	if ok {
		j.privateKey = privateKey
	}

	return j
}
