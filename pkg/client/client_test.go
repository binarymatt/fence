package client

import (
	"context"
	"errors"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/mock"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestNew(t *testing.T) {
	mock := NewMockFenceState(t)
	c := New(mock)
	must.NotNil(t, c.state)
}
func TestIsAllowed(t *testing.T) {
	principal := &fencev1.UID{}
	resource := &fencev1.UID{}
	action := &fencev1.UID{}
	cases := []struct {
		name string
		err  error
	}{
		{
			name: "happy path",
		},
		{
			name: "oops",
			err:  errors.New("oops"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockState := NewMockFenceState(t)
			mockState.EXPECT().IsAllowed(context.Background(), principal, action, resource).Return(tc.err)
			c := New(mockState)
			err := c.IsAllowed(context.Background(), principal, action, resource)
			must.ErrorIs(t, err, tc.err)
		})
	}
}
func TestIsAllowedFromContext(t *testing.T) {
	resource := &fencev1.UID{}
	action := &fencev1.UID{}
	errFailedAuth := errors.New("failed auth")
	cases := []struct {
		name        string
		ctx         context.Context
		expectation func(*MockFenceState)
		err         error
	}{
		{
			name: "happy path",
			ctx:  ContextWithPrincipal(context.Background(), &fencev1.UID{Type: "user", Id: "bob"}),
			expectation: func(m *MockFenceState) {
				m.EXPECT().IsAllowed(mock.AnythingOfType("*context.valueCtx"), &fencev1.UID{Type: "user", Id: "bob"}, action, resource).Return(nil)
			},
		},
		{
			name:        "no principal in context",
			ctx:         context.Background(),
			err:         ErrInvalidPrincipal,
			expectation: func(*MockFenceState) {},
		},
		{
			name: "failed auth",
			ctx:  ContextWithPrincipal(context.Background(), &fencev1.UID{Type: "user", Id: "bob"}),
			err:  errFailedAuth,
			expectation: func(m *MockFenceState) {
				m.EXPECT().IsAllowed(mock.AnythingOfType("*context.valueCtx"), &fencev1.UID{Type: "user", Id: "bob"}, resource, action).Return(errFailedAuth)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockState := NewMockFenceState(t)
			tc.expectation(mockState)
			c := New(mockState)
			err := c.IsAllowedFromContext(tc.ctx, resource, action)
			must.ErrorIs(t, err, tc.err)
		})
	}
}

func TestFenceToCedarUID(t *testing.T) {
	uid := fenceToCedarUID(&fencev1.UID{
		Type: "user",
		Id:   "bob",
	})
	must.Eq(t, cedar.NewEntityUID("user", "bob"), uid)
}
