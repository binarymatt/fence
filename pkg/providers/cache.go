package providers

import (
	"context"
	"sync"
	"time"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceProvider = (*CachedProvider)(nil)

type CachedProvider struct {
	fs              FenceProvider
	mu              sync.Mutex
	refreshDuration time.Duration
}

func (cfs *CachedProvider) IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) error {
	cfs.mu.Lock()
	defer cfs.mu.Unlock()
	return cfs.fs.IsAllowed(ctx, principal, action, resource)
}

func (cfs *CachedProvider) Refresh(ctx context.Context) error {
	cfs.mu.Lock()
	defer cfs.mu.Unlock()
	return cfs.fs.Refresh(ctx)
}

func NewCachedProvider(provider FenceProvider, refreshDuration time.Duration) (*CachedProvider, error) {
	wrapper := &CachedProvider{fs: provider, refreshDuration: refreshDuration}
	err := wrapper.Refresh(context.Background())
	if err != nil {
		return nil, err
	}
	return wrapper, nil
}
