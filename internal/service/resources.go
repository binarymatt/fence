package service

import (
	"context"
	"log/slog"

	"connectrpc.com/connect/v2"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func (s *Service) CreateResource(ctx context.Context, req *fencev1.CreateResourceRequest) (*fencev1.CreateResourceResponse, error) {
	if err := s.ds.AddEntities(ctx, req.GetResources()...); err != nil {
		slog.Error("failed to add resources", "error", err)
		return nil, connect.NewError(connect.CodeInternal, "failed to add resources")
	}
	return nil, nil
}
