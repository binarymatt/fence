package state

import (
	"context"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
)

var _ FenceState = (*RemoteServerState)(nil)

type RemoteServerConfig struct {
	Address string
	Timeout time.Duration
}

func NewRemoteServerState(cfg RemoteServerConfig) *RemoteServerState {
	retyrable := retryablehttp.NewClient()
	httpClient := retyrable.StandardClient()
	cl := fencev1connect.NewFenceServiceClient(httpClient, cfg.Address)
	return &RemoteServerState{client: cl}
}

type RemoteServerState struct {
	client fencev1connect.FenceServiceClient
}

func (a *RemoteServerState) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID) error {
	req := &fencev1.IsAllowedRequest{
		Principal: principal,
		Action:    action,
		Resource:  resource,
	}
	_, err := a.client.IsAllowed(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (a *RemoteServerState) Refresh(_ context.Context) error {
	return nil
}
func (a *RemoteServerState) refresh() error {
	return nil
}
