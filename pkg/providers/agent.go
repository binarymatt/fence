package providers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/scalalang2/golang-fifo/sieve"
	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
)

var _ FenceProvider = (*RemoteServerProvider)(nil)

type authorizedRoundTripper struct {
	baseTripper *retryablehttp.RoundTripper
	token       string
}

func (art *authorizedRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", art.token))
	res, err = art.baseTripper.RoundTrip(req)
	return
}

type RemoteServerConfig struct {
	Address     string
	BearerToken string
	Timeout     time.Duration
	CacheTTL    time.Duration
}

func NewRemoteServerProvider(cfg RemoteServerConfig) *RemoteServerProvider {
	retryable := retryablehttp.NewClient()
	httpClient := &http.Client{
		Transport: &authorizedRoundTripper{baseTripper: &retryablehttp.RoundTripper{Client: retryable}, token: cfg.BearerToken},
	}
	cl := fencev1connect.NewFenceServiceClient(httpClient, cfg.Address)
	return &RemoteServerProvider{client: cl}
}

type RemoteServerProvider struct {
	client fencev1connect.FenceServiceClient
}

func (a *RemoteServerProvider) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID, record map[string]*structpb.Value) (*fencev1.IsAllowedResponse, error) {
	req := &fencev1.IsAllowedRequest{
		Principal: principal,
		Action:    action,
		Resource:  resource,
		Context:   record,
	}
	resp, err := a.client.IsAllowed(ctx, req)
	if err != nil {
		return resp, NewAuthzError(principal, action, resource, err)
	}
	return resp, nil
}

func (a *RemoteServerProvider) Refresh(_ context.Context) error {
	return nil
}
