package service

import (
	"context"
	"slices"

	"buf.build/go/protovalidate"
	"connectrpc.com/connect/v2"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func (s *Service) CreateUser(ctx context.Context, req *fencev1.CreateUserRequest) (*fencev1.CreateUserResponse, error) {
	entity := req.User
	if req.GetTenant() != nil && !slices.Contains(entity.Parents, req.GetTenant()) {
		entity.Parents = append(entity.Parents, req.Tenant)
	}
	if err := s.ds.AddEntities(ctx, req.Tenant, entity); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err.Error())
	}
	return &fencev1.CreateUserResponse{}, nil
}
func (s *Service) DeleteUser(ctx context.Context, req *fencev1.DeleteUserRequest) (*fencev1.DeleteUserResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err.Error())
	}
	if err := s.ds.DeleteEntity(ctx, req.Tenant, req.Uid); err != nil {

		return nil, connect.NewError(connect.CodeInternal, err.Error())
	}
	return &fencev1.DeleteUserResponse{}, nil
}
func (s *Service) ListUsers(ctx context.Context, req *fencev1.ListUsersRequest) (*fencev1.ListUsersResponse, error) {
	return nil, nil
}
