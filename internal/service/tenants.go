package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect/v2"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

// TODO: add authn with admin token
func (s *Service) CreateTenant(ctx context.Context, req *fencev1.CreateTenantRequest) (*fencev1.CreateTenantResponse, error) {
	if err := s.ds.AddEntities(ctx, req.Tenant); err != nil {
		slog.Error("failed to add entities", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "failed to add entities")
	}
	return &fencev1.CreateTenantResponse{}, nil
}
func (s *Service) GetTenant(ctx context.Context, req *fencev1.GetTenantRequest) (*fencev1.GetTenantResponse, error) {
	return nil, nil
}
