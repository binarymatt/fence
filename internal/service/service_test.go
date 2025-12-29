package service

import (
	"context"
	"database/sql"
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

func TestParseUIDString(t *testing.T) {
	uid := parseUIDString(`User::"bob"`)
	expectedUID := cedar.NewEntityUID(cedar.EntityType("User"), cedar.String("bob"))
	must.Eq(t, expectedUID, uid)
}
func setupDB(t *testing.T, loadFixtures bool) *bun.DB {

	t.Helper()
	ctx := context.Background()
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
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
func setupTest(t *testing.T, loadFixtures bool) (*Service, *MockDataStore) {
	t.Helper()
	mockStore := NewMockDataStore(t)
	return &Service{ds: mockStore}, mockStore

}

func TestIsAllowed_Deny(t *testing.T) {
	s, _ := setupTest(t, true)
	req := &fencev1.IsAllowedRequest{
		Principal: &fencev1.UID{Id: "bob", Type: "User"},
		Action:    &fencev1.UID{Type: "Action", Id: "view"},
		Resource:  &fencev1.UID{Type: "Photo", Id: "VacationPhoto94.jpg"},
	}
	resp, err := s.IsAllowed(context.Background(), req)
	must.NoError(t, err)
	must.False(t, resp.Decision)
}
func TestIsAllowed_Allow(t *testing.T) {
	s, _ := setupTest(t, true)
	req := &fencev1.IsAllowedRequest{
		Principal: &fencev1.UID{Id: "alice", Type: "User"},
		Action:    &fencev1.UID{Type: "Action", Id: "view"},
		Resource:  &fencev1.UID{Type: "Photo", Id: "VacationPhoto94.jpg"},
	}
	resp, err := s.IsAllowed(context.Background(), req)
	must.NoError(t, err)
	must.True(t, resp.Decision)
}
