package providers

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceProvider = (*TestingProvider)(nil)

type TestingProvider struct {
	AllowCall bool
}

func (ts *TestingProvider) IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) (*fencev1.IsAllowedResponse, error) {
	return &fencev1.IsAllowedResponse{Decision: ts.AllowCall}, nil
}
func (ts *TestingProvider) Refresh(context.Context) error {
	return nil
}
