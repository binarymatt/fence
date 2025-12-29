package service

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var (
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

func (s *Service) CreateEntities(ctx context.Context, req *fencev1.CreateEntitiesRequest) (*fencev1.CreateEntitiesResponse, error) {
	for _, entity := range req.Entities {
		slog.Debug("creating entity", "uid", entity.Uid, "attrs", entity.GetAttributes())
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
		if err := s.ds.addEntity(ctx, entity); err != nil {
			slog.Error("failed to add entity", "record", dbEnt, "error", err)
			if errors.Is(err, ErrEntityAlreadyExists) {
				return nil, connect.NewError(connect.CodeInvalidArgument, err)
			}
			return nil, connect.NewError(connect.CodeUnknown, err)
		}
	}
	return &fencev1.CreateEntitiesResponse{}, nil
}
func (s *Service) DeleteEntity(ctx context.Context, req *fencev1.DeleteEntityRequest) (*fencev1.DeleteEntityResponse, error) {
	if err := s.ds.deleteEntity(ctx, req.GetUid()); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return &fencev1.DeleteEntityResponse{}, nil
}

func (s *Service) ListEntities(ctx context.Context, req *fencev1.ListEntitiesRequest) (*fencev1.ListEntitiesResponse, error) {
	entities, err := s.ds.getEntities(ctx)
	if err != nil {
		return nil, err
	}
	return &fencev1.ListEntitiesResponse{Entities: entities}, nil
}
func (s *Service) GetEntity(context.Context, *fencev1.GetEntityRequest) (*fencev1.GetEntityResponse, error) {
	return nil, nil
}
