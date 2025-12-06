package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var (
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

func (s *Service) CreateEntity(ctx context.Context, req *fencev1.CreateEntityRequest) (*fencev1.CreateEntityResponse, error) {
	parents := make([]UID, len(req.Entity.Parents))
	for i, ui := range req.Entity.Parents {
		parents[i] = fenceToDBUID(ui)
	}
	dbEnt := &Entity{
		Type:       req.GetEntity().GetUid().GetType(),
		ID:         req.GetEntity().GetUid().GetId(),
		Parents:    parents,
		Attributes: fenceToRecord(req.Entity.GetAttributes()),
		Tags:       fenceToRecord(req.Entity.Tags),
	}
	if err := s.addEntity(ctx, dbEnt); err != nil {
		if errors.Is(err, ErrEntityAlreadyExists) {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	return &fencev1.CreateEntityResponse{}, nil
}
func (s *Service) DeleteEntity(ctx context.Context, req *fencev1.DeleteEntityRequest) (*fencev1.DeleteEntityResponse, error) {
	if err := s.deleteEntity(ctx, req.GetUid().GetType(), req.GetUid().GetId()); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return &fencev1.DeleteEntityResponse{}, nil
}
