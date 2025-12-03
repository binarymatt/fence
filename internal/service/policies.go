package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func (s *service) CreatePolicy(ctx context.Context, req *fencev1.CreatePolicyRequest) (*fencev1.CreatePolicyResponse, error) {
	if err := s.addPolicy(ctx, req.Policy.GetId(), req.Policy.GetDefinition()); err != nil {
		if errors.Is(err, ErrPolicyAlreadyExists) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	return &fencev1.CreatePolicyResponse{}, nil
}
func (s *service) DeletePolicy(ctx context.Context, req *fencev1.DeletePolicyRequest) (*fencev1.DeletePolicyResponse, error) {
	if err := s.deletePolicy(ctx, req.Id); err != nil {
		if errors.Is(err, ErrPolicyNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	return &fencev1.DeletePolicyResponse{}, nil
}
