package service

import (
	"context"
	"database/sql"
	"errors"

	"connectrpc.com/connect"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func (s *Service) CreatePolicies(ctx context.Context, req *fencev1.CreatePoliciesRequest) (*fencev1.CreatePoliciesResponse, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(req.GetPolicies()))
	for i, policy := range req.GetPolicies() {

		if err := s.addPolicy(ctx, tx, policy.GetId(), policy.GetDefinition()); err != nil {
			if errors.Is(err, ErrPolicyAlreadyExists) {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}
			return nil, connect.NewError(connect.CodeUnknown, err)
		}
		ids[i] = policy.GetId()
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &fencev1.CreatePoliciesResponse{
		Ids: ids,
	}, nil
}
func (s *Service) DeletePolicy(ctx context.Context, req *fencev1.DeletePolicyRequest) (*fencev1.DeletePolicyResponse, error) {
	if err := s.deletePolicy(ctx, req.Id); err != nil {
		if errors.Is(err, ErrPolicyNotFound) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	return &fencev1.DeletePolicyResponse{}, nil
}

func (s *Service) ListPolicies(ctx context.Context, _ *fencev1.ListPoliciesRequest) (*fencev1.ListPoliciesResponse, error) {
	policies, err := s.getPolicies(ctx)
	if err != nil {
		return nil, err
	}
	protoPolicies := make([]*fencev1.Policy, len(policies))
	for i, p := range policies {
		protoPolicies[i] = p.ToProto()
	}

	return &fencev1.ListPoliciesResponse{
		Policies: protoPolicies,
	}, nil
}
func (s *Service) GetPolicy(context.Context, *fencev1.GetPolicyRequest) (*fencev1.GetPolicyResponse, error) {
	return nil, nil
}
