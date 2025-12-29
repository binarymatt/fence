package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestCreatePolicy(t *testing.T) {
	cases := []struct {
		name         string
		policy       *fencev1.Policy
		loadFixtures bool
		setupMock    func(*MockDataStore)
		err          error
	}{
		{
			name: "happy path",
			policy: &fencev1.Policy{
				Id:         "testPolicy",
				Definition: `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`,
			},
			setupMock: func(ds *MockDataStore) {
				ds.EXPECT().addPolicy(context.Background(), "testPolicy", `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`).Return(nil).Once()
			},
		},
		{
			name:         "existing policy",
			loadFixtures: true,
			policy: &fencev1.Policy{
				Id:         "policy0",
				Definition: `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`,
			},
			validate: func(t *testing.T, db *sql.DB, err error) {
				must.ErrorIs(t, err, ErrPolicyAlreadyExists)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, ds := setupTest(t, tc.loadFixtures)
			req := &fencev1.CreatePoliciesRequest{
				Policies: []*fencev1.Policy{tc.policy},
			}
			tc.setupMock(ds)
			_, err := s.CreatePolicies(context.Background(), req)
			must.ErrorIs(t, err, tc.err)
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		validate func(t *testing.T, db *sql.DB, err error)
	}{
		{
			name: "happy path",
			id:   "policy0",
			validate: func(t *testing.T, db *sql.DB, err error) {
				must.NoError(t, err)
				var count int
				db.QueryRow("select count(*) from policies").Scan(&count)
				must.Eq(t, 0, count)

			},
		},
		{
			name: "does not exist",
			id:   "testPolicy",
			validate: func(t *testing.T, db *sql.DB, err error) {
				must.ErrorIs(t, err, ErrPolicyNotFound)
				var count int
				db.QueryRow("select count(*) from policies").Scan(&count)
				must.Eq(t, 1, count)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			s, db := setupTest(t, true)

			req := &fencev1.DeletePolicyRequest{
				Id: tc.id,
			}
			_, err := s.DeletePolicy(context.Background(), req)
			tc.validate(t, db, err)
		})
	}

}
