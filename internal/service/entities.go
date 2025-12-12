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

func (s *Service) CreateEntities(ctx context.Context, req *fencev1.CreateEntitiesRequest) (*fencev1.CreateEntitiesResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	for _, entity := range req.Entities {
		parents := make([]UID, len(entity.Parents))
		for i, ui := range entity.Parents {
			parents[i] = fenceToDBUID(ui)
		}
		dbEnt := &Entity{
			Type:       entity.GetUid().GetType(),
			ID:         entity.GetUid().GetId(),
			Parents:    parents,
			Attributes: fenceToRecord(entity.GetAttributes()),
			Tags:       fenceToRecord(entity.Tags),
		}
		if err := s.addEntity(ctx, tx, dbEnt); err != nil {
			if errors.Is(err, ErrEntityAlreadyExists) {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}
			return nil, connect.NewError(connect.CodeUnknown, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &fencev1.CreateEntitiesResponse{}, nil
}
func (s *Service) DeleteEntity(ctx context.Context, req *fencev1.DeleteEntityRequest) (*fencev1.DeleteEntityResponse, error) {
	if err := s.deleteEntity(ctx, req.GetUid().GetType(), req.GetUid().GetId()); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return &fencev1.DeleteEntityResponse{}, nil
}

func (s *Service) ListEntities(ctx context.Context, req *fencev1.ListEntitiesRequest) (*fencev1.ListEntitiesResponse, error) {
	return nil, nil
}
func (s *Service) GetEntity(context.Context, *fencev1.GetEntityRequest) (*fencev1.GetEntityResponse, error) {
	return nil, nil
}
