package client

import (
	"context"
	"sync"
	"time"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceState = (*CachedState)(nil)

type CachedState struct {
	fs              FenceState
	mu              sync.Mutex
	refreshDuration time.Duration
}

func (cfs *CachedState) IsAllowed(ctx context.Context, principal, action, resource *fencev1.UID) error {
	cfs.mu.Lock()
	defer cfs.mu.Unlock()
	return cfs.fs.IsAllowed(ctx, principal, action, resource)
}

func (cfs *CachedState) Refresh(ctx context.Context) error {
	ticker := time.NewTicker(cfs.refreshDuration)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := cfs.refresh(); err != nil {
				return err
			}
		}
	}
}

func (cfs *CachedState) refresh() error {
	cfs.mu.Lock()
	defer cfs.mu.Unlock()
	return cfs.fs.refresh()
}

func NewCachedState(state FenceState, refreshDuration time.Duration) (*CachedState, error) {
	wrapper := &CachedState{fs: state, refreshDuration: refreshDuration}
	err := wrapper.refresh()
	if err != nil {
		return nil, err
	}
	return wrapper, nil
}
