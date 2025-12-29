package state

import (
	"context"
	"errors"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/pkg/providers"
)

func TestNewClient(t *testing.T) {
	mock := providers.NewMockFenceProvider(t)
	c := NewClient(mock)
	must.NotNil(t, c.provider)
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
			mockState := providers.NewMockFenceProvider(t)
			mockState.EXPECT().IsAllowed(context.Background(), principal, action, resource).Return(tc.err)
			c := NewClient(mockState)
			err := c.IsAllowed(context.Background(), principal, action, resource)
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
