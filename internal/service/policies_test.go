package service

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestCreatePolicy(t *testing.T) {
	cases := []struct {
		name         string
		policy       *fencev1.Policy
		loadFixtures bool
		setupMock    func(context.Context, *MockDataStore)
		err          error
	}{
		{
			name: "happy path",
			policy: &fencev1.Policy{
				Id:         "testPolicy",
				Definition: `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`,
			},
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				ds.EXPECT().addPolicy(ctx, "testPolicy", `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`).Return(nil).Once()
			},
		},
		{
			name:         "existing policy",
			loadFixtures: true,
			policy: &fencev1.Policy{
				Id:         "policy0",
				Definition: `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`,
			},
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				ds.EXPECT().addPolicy(ctx, "policy0", `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`).Return(ErrPolicyAlreadyExists)
			},
			err: ErrPolicyAlreadyExists,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, ds := setupTest(t)
			req := &fencev1.CreatePoliciesRequest{
				Policies: []*fencev1.Policy{tc.policy},
			}
			tc.setupMock(t.Context(), ds)
			_, err := s.CreatePolicies(t.Context(), req)
			must.ErrorIs(t, err, tc.err)
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	cases := []struct {
		name      string
		id        string
		setupMock func(context.Context, *MockDataStore)
		err       error
	}{
		{
			name: "happy path",
			id:   "policy0",
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				ds.EXPECT().deletePolicy(ctx, "policy0").Return(nil)
			},
		},
		{
			name: "does not exist",
			id:   "testPolicy",
			setupMock: func(ctx context.Context, ds *MockDataStore) {
				ds.EXPECT().deletePolicy(ctx, "testPolicy").Return(ErrPolicyNotFound)

			},
			err: ErrPolicyNotFound,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, ds := setupTest(t)
			tc.setupMock(t.Context(), ds)

			req := &fencev1.DeletePolicyRequest{
				Id: tc.id,
			}
			_, err := s.DeletePolicy(t.Context(), req)
			must.ErrorIs(t, err, tc.err)
		})
	}

}
