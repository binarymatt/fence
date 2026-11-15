package translation

import (
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestTranslateFromIsAllowedResponse_RoundTrip(t *testing.T) {
	cases := map[string]struct {
		decision cedar.Decision
		diag     cedar.Diagnostic
	}{
		"allow with reason": {
			decision: cedar.Allow,
			diag: cedar.Diagnostic{
				Reasons: []cedar.DiagnosticReason{{
					PolicyID: "policy0",
					Position: cedar.Position{Filename: "policies.cedar", Offset: 10, Line: 2, Column: 3},
				}},
			},
		},
		"deny from forbid with error": {
			decision: cedar.Deny,
			diag: cedar.Diagnostic{
				Reasons: []cedar.DiagnosticReason{{PolicyID: "forbid0"}},
				Errors: []cedar.DiagnosticError{{
					PolicyID: "policy1",
					Position: cedar.Position{Line: 4, Column: 1},
					Message:  "type error",
				}},
			},
		},
		"default deny": {
			decision: cedar.Deny,
			diag:     cedar.Diagnostic{},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			resp := TranslateAuthorizeResponse(tc.decision, tc.diag)
			decision, diag := TranslateFromIsAllowedResponse(resp)
			must.Eq(t, tc.decision, decision)
			must.Eq(t, tc.diag, diag)
		})
	}
}

func TestTranslateFromIsAllowedResponse_Nil(t *testing.T) {
	decision, diag := TranslateFromIsAllowedResponse(nil)
	must.Eq(t, cedar.Deny, decision)
	must.Eq(t, cedar.Diagnostic{}, diag)

	decision, diag = TranslateFromIsAllowedResponse(&fencev1.IsAllowedResponse{Decision: true})
	must.Eq(t, cedar.Allow, decision)
	must.Eq(t, cedar.Diagnostic{}, diag)
}

func TestMapToRecordRoundTrip(t *testing.T) {
	decimal, err := cedar.NewDecimalFromFloat(1.5)
	must.NoError(t, err)

	cases := map[string]cedar.Record{
		"empty": cedar.NewRecord(cedar.RecordMap{}),
		"primitives": cedar.NewRecord(cedar.RecordMap{
			"name":  cedar.String("alice"),
			"admin": cedar.True,
			"age":   cedar.Long(30),
		}),
		"nested": cedar.NewRecord(cedar.RecordMap{
			"tags":   cedar.NewSet(cedar.String("a"), cedar.String("b")),
			"nested": cedar.NewRecord(cedar.RecordMap{"level": cedar.Long(2)}),
		}),
		"entity and extension": cedar.NewRecord(cedar.RecordMap{
			"owner": cedar.NewEntityUID("User", "bob"),
			"ratio": decimal,
		}),
	}
	for name, record := range cases {
		t.Run(name, func(t *testing.T) {
			got := MapToRecord(RecordToMap(record))
			must.True(t, record.Equal(got), must.Sprintf("expected %s, got %s", record, got))
		})
	}
}
