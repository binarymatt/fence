package service

import (
	"context"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

// TODO: add authn with admin token
func (s *Service) CreateTenant(ctx context.Context, req *fencev1.CreateTenantRequest) (*fencev1.CreateTenantResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := s.ds.AddEntity(ctx, req.Tenant); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return &fencev1.CreateTenantResponse{}, nil
}
func (s *Service) GetTenant(ctx context.Context, req *fencev1.GetTenantRequest) (*fencev1.GetTenantResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return nil, nil
}
