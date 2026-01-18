package providers

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type FenceProvider interface {
	IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) (*fencev1.IsAllowedResponse, error)
	Refresh(context.Context) error
}
