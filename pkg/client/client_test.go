package client

import (
	"context"
	"errors"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/pkg/providers"
)

func TestNew(t *testing.T) {
	mock := providers.NewMockFenceProvider(t)
	c := New(mock)
	must.NotNil(t, c.provider)
}

func TestAuthorize(t *testing.T) {
	req := cedar.Request{
		Principal: cedar.NewEntityUID("User", "alice"),
		Action:    cedar.NewEntityUID("Action", "view"),
		Resource:  cedar.NewEntityUID("Photo", "VacationPhoto94.jpg"),
	}
	principal := NewUID("User", "alice")
	action := NewUID("Action", "view")
	resource := NewUID("Photo", "VacationPhoto94.jpg")

	cases := []struct {
		name         string
		resp         *fencev1.IsAllowedResponse
		err          error
		expected     cedar.Decision
		expectedDiag cedar.Diagnostic
	}{
		{
			name: "allow",
			resp: &fencev1.IsAllowedResponse{
				Decision: true,
				Diagnostics: &fencev1.Diagnostics{
					Reasons: []*fencev1.Reason{{PolicyId: "policy0", Position: &fencev1.Position{Line: 1, Column: 1}}},
				},
			},
			expected: cedar.Allow,
			expectedDiag: cedar.Diagnostic{
				Reasons: []cedar.DiagnosticReason{{PolicyID: "policy0", Position: cedar.Position{Line: 1, Column: 1}}},
			},
		},
		{
			name: "default deny",
			resp: &fencev1.IsAllowedResponse{
				Diagnostics: &fencev1.Diagnostics{
					Reasons: []*fencev1.Reason{{PolicyId: "no policy", Message: "default deny"}},
				},
			},
			expected: cedar.Deny,
		},
		{
			name:     "provider error fails closed",
			err:      errors.New("connection refused"),
			expected: cedar.Deny,
			expectedDiag: cedar.Diagnostic{
				Errors: []cedar.DiagnosticError{{Message: "connection refused"}},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockState := providers.NewMockFenceProvider(t)
			mockState.EXPECT().IsAllowed(context.Background(), principal, action, resource, map[string]*structpb.Value{}, []*fencev1.Entity{}).Return(tc.resp, tc.err)
			c := New(mockState)
			decision, diag := c.Authorize(req, nil)
			must.Eq(t, tc.expected, decision)
			must.Eq(t, tc.expectedDiag, diag)
		})
	}
}

func TestCedarToFenceUID(t *testing.T) {
	must.Eq(t, NewUID("User", "bob"), cedarToFenceUID(cedar.NewEntityUID("User", "bob")))
}

func TestFenceToCedarUID(t *testing.T) {
	uid := fenceToCedarUID(&fencev1.UID{
		Type: "user",
		Id:   "bob",
	})
	must.Eq(t, cedar.NewEntityUID("user", "bob"), uid)
}
