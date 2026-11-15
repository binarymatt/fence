package providers

import (
	"context"

	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type FenceProvider interface {
	IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID, record map[string]*structpb.Value, entities []*fencev1.Entity) (*fencev1.IsAllowedResponse, error)
	Refresh(context.Context) error
}
