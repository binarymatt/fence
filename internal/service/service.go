package service

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect/v2"
	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
	"github.com/binarymatt/fence/internal/database"
	"github.com/binarymatt/fence/internal/translation"
)

var _ fencev1connect.FenceServiceHandler = (*Service)(nil)
var _ fencev1connect.FenceAdminServiceHandler = (*Service)(nil)

func New(ds database.DataStore) *Service {
	return &Service{ds: ds}
}

type Service struct {
	ds database.DataStore
}

func fenceToCedarUID(uid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(uid.GetType()), cedar.String(uid.GetId()))
}
func (s *Service) addEntities(em cedar.EntityMap, entities []*fencev1.Entity) {
	for _, e := range entities {
		id := cedar.NewEntityUID(cedar.EntityType(e.Uid.Type), cedar.String(e.Uid.Id))

		parentRecords := make([]cedar.EntityUID, len(e.Parents))
		for i, uid := range e.Parents {
			parentRecords[i] = cedar.NewEntityUID(cedar.EntityType(uid.Type), cedar.String(uid.Id))
		}

		parents := cedar.NewEntityUIDSet(parentRecords...)
		ent := cedar.Entity{
			UID:        id,
			Parents:    parents,
			Attributes: translation.MapToRecord(e.Attributes),
			Tags:       translation.MapToRecord(e.Tags),
		}
		em[id] = ent
	}
}
func (s *Service) IsAllowed(ctx context.Context, req *fencev1.IsAllowedRequest) (*fencev1.IsAllowedResponse, error) {
	tenant := &fencev1.UID{
		Type: "Tenant",
		Id:   req.GetTenant(),
	}
	ps, err := s.ds.GetPolicySet(ctx, tenant)
	if err != nil {
		slog.Error("failed to get policy set", "eror", err)
		return nil, connect.NewError(connect.CodeInternal, "failed to get policy set")
	}
	em, err := s.ds.GetEntityMap(ctx, tenant)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, "failed to get entity map")
	}
	// Add additional passed in entities
	s.addEntities(em, req.GetEntities())

	principal := req.GetPrincipal()
	action := req.GetAction()
	resource := req.GetResource()

	cedarContext := translation.MapToRecord(req.GetContext())

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
