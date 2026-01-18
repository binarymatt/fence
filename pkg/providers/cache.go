package providers

import (
	"context"
	"fmt"
	"sync"
	"time"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceProvider = (*CachedProvider)(nil)

type CachedProvider struct {
	ip              FenceProvider
	mu              sync.Mutex
	refreshDuration time.Duration
}

func (cfs *CachedProvider) IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) (*fencev1.IsAllowedResponse, error) {
	cfs.mu.Lock()
	defer cfs.mu.Unlock()
	return cfs.ip.IsAllowed(ctx, principal, action, resource)
}

func (cfs *CachedProvider) Refresh(ctx context.Context) error {
	ticker := time.NewTicker(cfs.refreshDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done with context")
			return nil
		case <-ticker.C:
			if err := cfs.refreshInternalProvider(ctx); err != nil {
				return err
			}
		}
	}
}
func (cfs *CachedProvider) refreshInternalProvider(ctx context.Context) error {
	cfs.mu.Lock()
	defer cfs.mu.Unlock()
	return cfs.ip.Refresh(ctx)
}

func NewCachedProvider(provider FenceProvider, refreshDuration time.Duration) (*CachedProvider, error) {
	wrapper := &CachedProvider{ip: provider, refreshDuration: refreshDuration}
	err := wrapper.refreshInternalProvider(context.Background())
	if err != nil {
		return nil, err
	}
	return wrapper, nil
}
