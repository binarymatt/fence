package service

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
	"github.com/binarymatt/fence/internal/translation"
)

var _ fencev1connect.FenceServiceHandler = (*Service)(nil)
var _ fencev1connect.FenceAdminServiceHandler = (*Service)(nil)

func New(ds DataStore) *Service {
	return &Service{ds: ds}
}

type Service struct {
	ds DataStore
}

func fenceToCedarUID(uid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(uid.GetType()), cedar.String(uid.GetId()))
}
func (s *Service) IsAllowed(ctx context.Context, req *fencev1.IsAllowedRequest) (*fencev1.IsAllowedResponse, error) {
	ps, err := s.ds.getPolicySet(ctx)
	if err != nil {
		slog.Error("failed to get policy set", "eror", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	em, err := s.ds.getEntityMap(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	principal := req.GetPrincipal()
	action := req.GetAction()
	resource := req.GetResource()

	cedarContext := fenceToRecord(req.GetContext())

	cedarReq := cedar.Request{
		Principal: fenceToCedarUID(principal),
		Action:    fenceToCedarUID(action),
		Resource:  fenceToCedarUID(resource),
		Context:   cedarContext,
	}
	decision, diag := cedar.Authorize(ps, em, cedarReq)
	slog.Debug("authorize call finished", "decision", decision, "diagnostics", diag, "request", req)
	resp := translation.TranslateAuthorizeResponse(decision, diag)
	return resp, nil
}
func parseUIDString(uidStr string) cedar.EntityUID {
	parts := strings.Split(uidStr, "::")
	uidType := parts[0]
	id := parts[1]
	if strings.HasPrefix(id, "\"") && strings.HasSuffix(id, "\"") {
		id = id[1 : len(id)-1]
	}
	return cedar.NewEntityUID(cedar.EntityType(uidType), cedar.String(id))
}
