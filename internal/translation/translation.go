package translation

import (
	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TranslateAuthorizeResponse(decision cedar.Decision, diag cedar.Diagnostic) *fencev1.IsAllowedResponse {
	reasons := make([]*fencev1.Reason, len(diag.Reasons))
	for i, r := range diag.Reasons {
		reasons[i] = cedarToFenceReason(r)
	}
	if len(reasons) == 0 && decision == false {
		reasons = append(reasons, &fencev1.Reason{PolicyId: "no policy", Message: "default deny"})
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
