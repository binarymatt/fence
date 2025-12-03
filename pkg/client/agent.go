package client

import (
	"context"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
)

var _ FenceState = (*FenceAgentState)(nil)

type AgentConfig struct {
	Address string
	Timeout time.Duration
}

func NewAgentState(cfg AgentConfig) *FenceAgentState {
	retyrable := retryablehttp.NewClient()
	httpClient := retyrable.StandardClient()
	cl := fencev1connect.NewFenceServiceClient(httpClient, cfg.Address)
	return &FenceAgentState{client: cl}
}

type FenceAgentState struct {
	client fencev1connect.FenceServiceClient
}

func (a *FenceAgentState) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID) error {
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

func (a *FenceAgentState) Refresh(_ context.Context) error {
	return nil
}
func (a *FenceAgentState) refresh() error {
	return nil
}
