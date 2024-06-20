package authn

import (
	"context"
	"encoding/json"

	"github.com/w-h-a/pkg/utils/metadatautils"
)

const (
	accountKey = "auth-account"
)

func ContextWithAccount(ctx context.Context, account *Account) (context.Context, error) {
	bytes, err := json.Marshal(account)
	if err != nil {
		return ctx, err
	}

	return metadatautils.SetContext(ctx, accountKey, string(bytes)), nil
}

func AccountFromContext(ctx context.Context) (*Account, error) {
	str, ok := metadatautils.GetContext(ctx, accountKey)
	if !ok {
		return nil, nil
	}

	acc := &Account{}

	if err := json.Unmarshal([]byte(str), acc); err != nil {
		return nil, err
	}

	return acc, nil
}
