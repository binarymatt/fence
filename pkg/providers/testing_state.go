package providers

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceProvider = (*TestingProvider)(nil)

type TestingProvider struct {
	AllowCall bool
}

func (ts *TestingProvider) IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) error {
	if !ts.AllowCall {
		return NewAuthzError(principal, action, resource)
	}
	return nil
}
func (ts *TestingProvider) Refresh(context.Context) error {
	return nil
}
