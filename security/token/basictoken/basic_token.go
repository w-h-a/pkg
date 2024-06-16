package basictoken

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/w-h-a/pkg/security/token"
	"github.com/w-h-a/pkg/store"
)

type basicTokenProvider struct {
	options token.TokenOptions
	store   store.Store
}

func (t *basicTokenProvider) Options() token.TokenOptions {
	return t.options
}

func (t *basicTokenProvider) Generate(opts ...token.GenerateOption) (*token.Token, error) {
	// get the options
	options := token.NewGenerateOptions(opts...)

	// construct token
	tk := token.Token{
		AccessToken: uuid.New().String(),
		Created:     time.Now(),
		Expiry:      time.Now().Add(options.Expiry),
		Id:          options.Id,
		Roles:       options.Roles,
		Metadata:    options.Metadata,
	}

	// marshal the token
	bytes, err := json.Marshal(tk)
	if err != nil {
		return nil, err
	}

	// write to the token store
	if err := t.store.Write(&store.Record{
		Key:    tk.AccessToken,
		Value:  bytes,
		Expiry: options.Expiry,
	}); err != nil {
		return nil, err
	}

	return &tk, nil
}

func (t *basicTokenProvider) Inspect(tok string) (*token.Token, error) {
	// lookup the token
	records, err := t.store.Read(tok)
	if err == store.ErrRecordNotFound {
		return nil, token.ErrInvalidToken
	} else if err != nil {
		return nil, err
	}

	// unmarshal bytes
	bytes := records[0].Value
	var tk *token.Token
	if err := json.Unmarshal(bytes, &tk); err != nil {
		return nil, err
	}

	// ensure the token hasn't expired
	// this should be handled by the store, but let's check again
	if tk.Expiry.Unix() < time.Now().Unix() {
		return nil, token.ErrInvalidToken
	}

	return tk, nil
}

func (t *basicTokenProvider) String() string {
	return "basic"
}

func NewTokenProvider(opts ...token.TokenOption) token.TokenProvider {
	options := token.NewTokenOptions(opts...)

	b := &basicTokenProvider{
		options: options,
	}

	s, ok := GetStoreFromContext(options.Context)
	if ok {
		b.store = s
	}

	return b
}
