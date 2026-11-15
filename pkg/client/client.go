package client

import (
	"context"
	"sync"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/internal/translation"
	"github.com/binarymatt/fence/pkg/providers"
)

type Client struct {
	mu       sync.Mutex
	provider providers.FenceProvider
}

// Authorize mirrors cedar.Authorize but evaluates the request through the
// client's provider. It fails closed: if the provider returns an error, the
// decision is Deny and the error is reported in the diagnostic.
//
// req.Context is not forwarded, because FenceProvider.IsAllowed does not yet
// accept a Cedar context.
func (c *Client) Authorize(req cedar.Request, entities []cedar.Entity) (cedar.Decision, cedar.Diagnostic) {
	fenceEntities := make([]*fencev1.Entity, len(entities))
	for i, ent := range entities {
		fenceEntities[i] = translation.CedarEntityToFence(ent)
	}
	resp, err := c.provider.IsAllowed(
		context.Background(),
		cedarToFenceUID(req.Principal),
		cedarToFenceUID(req.Action),
		cedarToFenceUID(req.Resource),
		translation.RecordToMap(req.Context),
		fenceEntities,
	)
	if err != nil {
		return cedar.Deny, cedar.Diagnostic{
			Errors: []cedar.DiagnosticError{{Message: err.Error()}},
		}
	}
	return translation.TranslateFromIsAllowedResponse(resp)
}

func New(provider providers.FenceProvider) *Client {
	return &Client{provider: provider}
}
func fenceToCedarUID(fuid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(fuid.GetType()), cedar.String(fuid.GetId()))
}
func cedarToFenceUID(uid cedar.EntityUID) *fencev1.UID {
	return NewUID(string(uid.Type), string(uid.ID))
}
