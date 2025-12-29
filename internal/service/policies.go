package service

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func (s *Service) CreatePolicies(ctx context.Context, req *fencev1.CreatePoliciesRequest) (*fencev1.CreatePoliciesResponse, error) {
	ids := make([]string, len(req.GetPolicies()))
	slog.Debug("creating policies", "data", req.Policies)

	for i, policy := range req.GetPolicies() {
		if err := s.ds.addPolicy(ctx, policy.GetId(), policy.GetDefinition()); err != nil {
			if errors.Is(err, ErrPolicyAlreadyExists) {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}
			return nil, connect.NewError(connect.CodeUnknown, err)
		}
		ids[i] = policy.GetId()
	}
	return &fencev1.CreatePoliciesResponse{
		Ids: ids,
	}, nil
}
func (s *Service) DeletePolicy(ctx context.Context, req *fencev1.DeletePolicyRequest) (*fencev1.DeletePolicyResponse, error) {
	if err := s.ds.deletePolicy(ctx, req.Id); err != nil {
		if errors.Is(err, ErrPolicyNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	return &fencev1.DeletePolicyResponse{}, nil
}

func (s *Service) ListPolicies(ctx context.Context, _ *fencev1.ListPoliciesRequest) (*fencev1.ListPoliciesResponse, error) {
	slog.Debug("listing policies")
	policies, err := s.ds.getPolicies(ctx)
	if err != nil {
		return nil, err
	}

	return &fencev1.ListPoliciesResponse{
		Policies: policies,
	}, nil
}
func (s *Service) GetPolicy(ctx context.Context, req *fencev1.GetPolicyRequest) (*fencev1.GetPolicyResponse, error) {
	policy, err := s.ds.getPolicy(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	data := policy.MarshalCedar()
	p := &fencev1.Policy{
		Id:         req.GetId(),
		Definition: string(data),
	}
	return &fencev1.GetPolicyResponse{Policy: p}, nil
}
