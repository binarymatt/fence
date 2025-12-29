package state

import (
	"context"
	"sync"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/pkg/providers"
)

type Client struct {
	mu       sync.Mutex
	provider providers.FenceProvider
}

func (c *Client) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID) error {
	return c.provider.IsAllowed(ctx, principal, resource, action)
}

func NewClient(provider providers.FenceProvider) *Client {
	return &Client{provider: provider}
}
func fenceToCedarUID(fuid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(fuid.GetType()), cedar.String(fuid.GetId()))
}
