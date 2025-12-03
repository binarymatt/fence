package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestCreatePolicy(t *testing.T) {
	cases := []struct {
		name         string
		policy       *fencev1.Policy
		validate     func(t *testing.T, db *sql.DB, err error)
		loadFixtures bool
	}{
		{
			name: "happy path",
			policy: &fencev1.Policy{
				Id:         "testPolicy",
				Definition: `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`,
			},
			validate: func(t *testing.T, db *sql.DB, err error) {
				must.NoError(t, err)
				var dbID, dbContent string
				n := time.Now().UTC()
				var createdAt, updatedAt time.Time
				db.QueryRow("Select id, content, created_at, updated_at from policies").Scan(&dbID, &dbContent, &createdAt, &updatedAt)
				must.Eq(t, "testPolicy", dbID)
				must.Eq(t, `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`, dbContent)
				must.Eq(t, createdAt.Format(time.RFC3339), n.Format(time.RFC3339))
				must.Eq(t, updatedAt.Format(time.RFC3339), n.Format(time.RFC3339))

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
			s, db := setupTest(t, tc.loadFixtures)
			req := &fencev1.CreatePolicyRequest{
				Policy: tc.policy,
			}
			_, err := s.CreatePolicy(context.Background(), req)
			tc.validate(t, db, err)
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
