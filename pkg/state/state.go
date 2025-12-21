package state

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type FenceState interface {
	IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) error
	Refresh(context.Context) error
	refresh() error
}

// TODO: CachedObjectStoreState - same as cached file, but from object store
// TODO: ServerState - all calls are proxies to server
// TODO: CachedServerState - background job regularly updates local state from server
