package service

import (
	"context"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
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
	cedarReq := cedar.Request{
		Principal: fenceToCedarUID(principal),
		Action:    fenceToCedarUID(action),
		Resource:  fenceToCedarUID(resource),
	}
	decision, diag := cedar.Authorize(ps, em, cedarReq)
	slog.Debug("authorize call finished", "decision", decision, "diagnostics", diag, "request", req)
	reasons := make([]*fencev1.Reason, len(diag.Reasons))
	for i, r := range diag.Reasons {
		reasons[i] = cedarToFenceReason(r)
	}
	errors := make([]*fencev1.Error, len(diag.Errors))
	for i, err := range diag.Errors {
		errors[i] = cedarToFenceError(err)
	}
	resp := &fencev1.IsAllowedResponse{
		Decision: bool(decision),
		Diagnostics: &fencev1.Diagnostics{
			Reasons: reasons,
			Errors:  errors,
		},
	}
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
func cedarToFenceReason(reason cedar.DiagnosticReason) *fencev1.Reason {
	return &fencev1.Reason{
		PolicyId: string(reason.PolicyID),
		Position: &fencev1.Position{
			FileName: reason.Position.Filename,
			Line:     int64(reason.Position.Line),
			Column:   int64(reason.Position.Column),
			Offset:   int64(reason.Position.Offset),
		},
	}
}
func cedarToFenceError(err cedar.DiagnosticError) *fencev1.Error {
	return &fencev1.Error{
		PolicyId: string(err.PolicyID),
		Position: &fencev1.Position{
			FileName: err.Position.Filename,
			Line:     int64(err.Position.Line),
			Column:   int64(err.Position.Column),
			Offset:   int64(err.Position.Offset),
		},
		Message: err.Message,
	}
}
