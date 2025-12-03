package client

import (
	"context"
	"sync"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type client struct {
	mu    sync.Mutex
	state FenceState
}

func (c *client) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID) error {
	return c.state.IsAllowed(ctx, principal, resource, action)
}

func (c *client) IsAllowedFromContext(ctx context.Context, action *fencev1.UID, resource *fencev1.UID) error {
	principal, err := PrincipalFromContext(ctx)
	if err != nil {
		return err
	}
	return c.IsAllowed(ctx, principal, action, resource)
}
func New(state FenceState) *client {
	return &client{state: state}
}
func fenceToCedarUID(fuid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(fuid.GetType()), cedar.String(fuid.GetId()))
}
