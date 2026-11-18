package database

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

	"github.com/binarymatt/fence/internal/models"
)

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
	_, err = db.NewCreateTable().Model((*models.Entity)(nil)).IfNotExists().Exec(ctx)
	_, err = db.NewCreateTable().Model((*models.Policy)(nil)).IfNotExists().Exec(ctx)
	if loadFixtures {
		fixture := dbfixture.New(db)

		fs := os.DirFS("../../")
		if err := fixture.Load(ctx, fs, "testdata/fixture.yaml"); err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func TestGetPolicySet(t *testing.T) {
	db := setupDB(t, true)
	defer db.Close()
	ds := NewSqlDatastore(db)
	ps, err := ds.GetPolicySet(context.Background())
	must.NoError(t, err)
	p := ps.Get("policy0")
	data := p.MarshalCedar()
	expectedPolicy := `permit (
    principal == User::"alice",
    action == Action::"view",
    resource in Album::"jane_vacation"
);`
	must.Eq(t, expectedPolicy, string(data))
}

func TestGetEntities(t *testing.T) {
	ctx := context.Background()
	db := setupDB(t, true)
	defer db.Close()
	ds := NewSqlDatastore(db)

	bob := cedar.Entity{
		UID:     cedar.NewEntityUID(cedar.EntityType("User"), cedar.String("bob")),
		Parents: cedar.NewEntityUIDSet(cedar.NewEntityUID(cedar.EntityType("Group"), cedar.String("people"))),
	}
	alice := cedar.Entity{
		UID: cedar.NewEntityUID(cedar.EntityType("User"), cedar.String("alice")),
	}
	photo := cedar.Entity{
		UID:     cedar.NewEntityUID(cedar.EntityType("Photo"), cedar.String("VacationPhoto94.jpg")),
		Parents: cedar.NewEntityUIDSet(cedar.NewEntityUID(cedar.EntityType("Album"), cedar.String("jane_vacation"))),
	}
	expectedMap := cedar.EntityMap{
		bob.UID:   bob,
		alice.UID: alice,
		photo.UID: photo,
	}
	em, err := ds.GetEntityMap(ctx)
	must.NoError(t, err)
	must.Eq(t, expectedMap, em)
}
