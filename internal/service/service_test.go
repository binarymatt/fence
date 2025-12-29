package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

const entitiesJSON = `[
  {
    "uid": { "type": "User", "id": "alice" },
    "attrs": { "age": 18 },
    "parents": []
  },
  {
    "uid": { "type": "Photo", "id": "VacationPhoto94.jpg" },
    "attrs": {},
    "parents": [{ "type": "Album", "id": "jane_vacation" }]
  },
  {
  	"uid": {"type":"User","id":"bob"},
  	"parents": [{"type":"Group","id":"people"}]
  }
]`
const policyCedar = `permit(principal == User::"alice",action == Action::"view",resource in Album::"jane_vacation");`

func policySet() *cedar.PolicySet {

	ps := cedar.NewPolicySet()
	var policy cedar.Policy
	policy.UnmarshalCedar([]byte(policyCedar))
	ps.Add("policy0", &policy)
	return ps
}
func entityMap() cedar.EntityMap {
	var entities cedar.EntityMap
	json.Unmarshal([]byte(entitiesJSON), &entities)
	return entities
}
func TestParseUIDString(t *testing.T) {
	uid := parseUIDString(`User::"bob"`)
	expectedUID := cedar.NewEntityUID(cedar.EntityType("User"), cedar.String("bob"))
	must.Eq(t, expectedUID, uid)
}
func setupDB(t *testing.T, loadFixtures bool) *bun.DB {

	t.Helper()
	ctx := context.Background()
	sqldb, err := sql.Open(sqliteshim.ShimName, ":memory:")
	// sqldb, err := sql.Open(sqliteshim.ShimName, "../../test.db")
	must.NoError(t, err)
	t.Cleanup(func() {
		sqldb.Close()
	})

	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(false),
	))
	_, err = db.NewCreateTable().Model((*Entity)(nil)).IfNotExists().Exec(ctx)
	_, err = db.NewCreateTable().Model((*Policy)(nil)).IfNotExists().Exec(ctx)
	if loadFixtures {
		fixture := dbfixture.New(db)

		fs := os.DirFS("../../")
		if err := fixture.Load(ctx, fs, "testdata/fixture.yaml"); err != nil {
			t.Fatal(err)
		}
	}
	return db
}
func setupTest(t *testing.T) (*Service, *MockDataStore) {
	t.Helper()
	mockStore := NewMockDataStore(t)
	return &Service{ds: mockStore}, mockStore

}

func TestIsAllowed_Deny(t *testing.T) {
	s, ds := setupTest(t)
	req := &fencev1.IsAllowedRequest{
		Principal: &fencev1.UID{Id: "bob", Type: "User"},
		Action:    &fencev1.UID{Type: "Action", Id: "view"},
		Resource:  &fencev1.UID{Type: "Photo", Id: "VacationPhoto94.jpg"},
	}
	ds.EXPECT().getPolicySet(t.Context()).Return(policySet(), nil)
	ds.EXPECT().getEntityMap(t.Context()).Return(entityMap(), nil)
	resp, err := s.IsAllowed(t.Context(), req)
	must.NoError(t, err)
	must.False(t, resp.Decision)
}
func TestIsAllowed_Allow(t *testing.T) {
	s, ds := setupTest(t)
	req := &fencev1.IsAllowedRequest{
		Principal: &fencev1.UID{Id: "alice", Type: "User"},
		Action:    &fencev1.UID{Type: "Action", Id: "view"},
		Resource:  &fencev1.UID{Type: "Photo", Id: "VacationPhoto94.jpg"},
	}
	ds.EXPECT().getPolicySet(t.Context()).Return(policySet(), nil)
	ds.EXPECT().getEntityMap(t.Context()).Return(entityMap(), nil)
	resp, err := s.IsAllowed(t.Context(), req)
	must.NoError(t, err)
	must.True(t, resp.Decision)
}
