package fence

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceState = (*TestingState)(nil)

type TestingState struct {
	AllowCall bool
}

func (ts *TestingState) IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) error {
	if !ts.AllowCall {
		return NewAuthzError(principal, action, resource)
	}
	return nil
}
func (ts *TestingState) Refresh(context.Context) error {
	return nil
}
func (ts *TestingState) refresh() error {
	return nil
}
