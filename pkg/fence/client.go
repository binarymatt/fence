package fence

import (
	"context"
	"sync"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type Client struct {
	mu    sync.Mutex
	state FenceState
}

func (c *Client) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID) error {
	return c.state.IsAllowed(ctx, principal, resource, action)
}

func (c *Client) IsAllowedFromContext(ctx context.Context, action *fencev1.UID, resource *fencev1.UID) error {
	principal, err := PrincipalFromContext(ctx)
	if err != nil {
		return err
	}
	return c.IsAllowed(ctx, principal, action, resource)
}
func NewClient(state FenceState) *Client {
	return &Client{state: state}
}
func fenceToCedarUID(fuid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(fuid.GetType()), cedar.String(fuid.GetId()))
}
