package translation

import (
	"encoding/json"
	"log/slog"

	"github.com/cedar-policy/cedar-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TranslateAuthorizeResponse(decision cedar.Decision, diag cedar.Diagnostic) *fencev1.IsAllowedResponse {
	reasons := make([]*fencev1.Reason, len(diag.Reasons))
	for i, r := range diag.Reasons {
		reasons[i] = cedarToFenceReason(r)
	}
	if len(reasons) == 0 && decision == false {
		reasons = append(reasons, defaultDenyReason())
	}
	errors := make([]*fencev1.Error, len(diag.Errors))
	for i, err := range diag.Errors {
		errors[i] = cedarToFenceError(err)
	}

	return &fencev1.IsAllowedResponse{
		Decision: bool(decision),
		Diagnostics: &fencev1.Diagnostics{
			Reasons: reasons,
			Errors:  errors,
		},
	}
}

// TranslateFromIsAllowedResponse is the inverse of TranslateAuthorizeResponse.
// The synthetic "default deny" reason added for denials with no matching
// policy is dropped so the result matches what cedar.Authorize would return.
// A nil response is treated as a deny with no diagnostics.
func TranslateFromIsAllowedResponse(resp *fencev1.IsAllowedResponse) (cedar.Decision, cedar.Diagnostic) {
	decision := cedar.Decision(resp.GetDecision())
	diag := cedar.Diagnostic{}

	for _, r := range resp.GetDiagnostics().GetReasons() {
		if isDefaultDenyReason(r) {
			continue
		}
		diag.Reasons = append(diag.Reasons, fenceToCedarReason(r))
	}
	for _, e := range resp.GetDiagnostics().GetErrors() {
		diag.Errors = append(diag.Errors, fenceToCedarError(e))
	}
	return decision, diag
}

const (
	defaultDenyPolicyID = "no policy"
	defaultDenyMessage  = "default deny"
)

func defaultDenyReason() *fencev1.Reason {
	return &fencev1.Reason{PolicyId: defaultDenyPolicyID, Message: defaultDenyMessage}
}

func isDefaultDenyReason(r *fencev1.Reason) bool {
	return r.GetPolicyId() == defaultDenyPolicyID && r.GetMessage() == defaultDenyMessage && r.GetPosition() == nil
}

func fenceToCedarPosition(p *fencev1.Position) cedar.Position {
	return cedar.Position{
		Filename: p.GetFileName(),
		Offset:   int(p.GetOffset()),
		Line:     int(p.GetLine()),
		Column:   int(p.GetColumn()),
	}
}

func fenceToCedarReason(r *fencev1.Reason) cedar.DiagnosticReason {
	return cedar.DiagnosticReason{
		PolicyID: cedar.PolicyID(r.GetPolicyId()),
		Position: fenceToCedarPosition(r.GetPosition()),
	}
}

func fenceToCedarError(e *fencev1.Error) cedar.DiagnosticError {
	return cedar.DiagnosticError{
		PolicyID: cedar.PolicyID(e.GetPolicyId()),
		Position: fenceToCedarPosition(e.GetPosition()),
		Message:  e.GetMessage(),
	}
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

// mapToRecord is the inverse of recordToMap: it encodes the fields as JSON and
// decodes them with Cedar's JSON format, so entity references ("__entity") and
// extension values ("__extn") round-trip.
func MapToRecord(values map[string]*structpb.Value) cedar.Record {
	data, err := protojson.Marshal(&structpb.Struct{Fields: values})
	if err != nil {
		slog.Error("could not marshal struct fields", "error", err)
		return cedar.Record{}
	}
	var r cedar.Record
	if err := r.UnmarshalJSON(data); err != nil {
		slog.Error("could not unmarshal cedar record", "error", err)
		return cedar.Record{}
	}
	return r
}
func RecordToMap(r cedar.Record) map[string]*structpb.Value {
	data, err := r.MarshalJSON()
	if err != nil {
		slog.Error("could not marshal cedar record")
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		slog.Error("could not unmarshal record data")
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		slog.Error("coudld not create struct")
	}
	return s.GetFields()
}
func FenceToCedarUID(uid *fencev1.UID) cedar.EntityUID {
	return cedar.NewEntityUID(cedar.EntityType(uid.GetType()), cedar.String(uid.GetId()))
}

func FenceToEntity(e *fencev1.Entity) cedar.Entity {
	uid := FenceToCedarUID(e.Uid)
	attrs := MapToRecord(e.GetAttributes())
	tags := MapToRecord(e.GetTags())

	parentSlice := make([]cedar.EntityUID, len(e.GetParents()))
	for i, e := range e.GetParents() {
		parentSlice[i] = FenceToCedarUID(e)
	}
	parents := cedar.NewEntityUIDSet(parentSlice...)
	return cedar.Entity{
		UID:        uid,
		Tags:       tags,
		Attributes: attrs,
		Parents:    parents,
	}
}
func CedarEntityToFence(c cedar.Entity) *fencev1.Entity {
	parents := c.Parents.Slice()
	fenceParents := make([]*fencev1.UID, len(parents))
	for i, uid := range parents {
		fenceParents[i] = CedarToFenceUID(uid)
	}
	attrs := RecordToMap(c.Attributes)
	tags := RecordToMap(c.Tags)
	return &fencev1.Entity{
		Uid:        CedarToFenceUID(c.UID),
		Parents:    fenceParents,
		Attributes: attrs,
		Tags:       tags,
	}
}

func newUID(uidType, id string) *fencev1.UID {
	return &fencev1.UID{
		Type: uidType,
		Id:   id,
	}
}
func CedarToFenceUID(uid cedar.EntityUID) *fencev1.UID {
	return newUID(string(uid.Type), string(uid.ID))
}
