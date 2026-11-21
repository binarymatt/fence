package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect/v2"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

// TODO: add authn with admin token
func (s *Service) CreateTenant(ctx context.Context, req *fencev1.CreateTenantRequest) (*fencev1.CreateTenantResponse, error) {
	if err := s.ds.AddTenant(ctx, req.Tenant); err != nil {
		slog.Error("failed to add entities", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "failed to add entities").WithCause(err)
	}
	return &fencev1.CreateTenantResponse{}, nil
}
func (s *Service) ListTenants(ctx context.Context, _ *fencev1.ListTenantRequest) (*fencev1.ListTenantResponse, error) {
	tenants, err := s.ds.ListTenants(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, "failed to get tenants").WithCause(err)
	}
	return &fencev1.ListTenantResponse{
		Tenants: tenants,
	}, nil
}
func (s *Service) GetTenant(ctx context.Context, req *fencev1.GetTenantRequest) (*fencev1.GetTenantResponse, error) {
	tenant, err := s.ds.GetTenant(ctx, req.Id)
	if err != nil {
		slog.Error("failed to retrieve tenant")
		return nil, connect.NewError(connect.CodeInternal, "failed to get tenant").WithCause(err)

	}
	return &fencev1.GetTenantResponse{Tenant: tenant}, nil
}
