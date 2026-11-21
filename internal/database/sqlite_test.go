package database

import (
	"context"
	"database/sql"
	"iter"
	"log/slog"
	"os"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/protobuf/testing/protocmp"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/internal/migrations"
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
	err = migrations.Migrate(ctx, db)
	must.NoError(t, err)
	// _, err = db.NewCreateTable().Model((*models.Entity)(nil)).IfNotExists().Exec(ctx)
	// _, err = db.NewCreateTable().Model((*models.Policy)(nil)).IfNotExists().Exec(ctx)
	if loadFixtures {
		db.RegisterModel((*models.Entity)(nil), (*models.Policy)(nil), (*models.Tenant)(nil))
		fixture := dbfixture.New(db)

		fs := os.DirFS("../../")
		if err := fixture.Load(ctx, fs, "testdata/fixture.yaml"); err != nil {
			slog.Error("failed to load fixtures", "error", err)
			t.Fatal(err)
		}
	}
	return db
}

func count(it iter.Seq2[cedar.PolicyID, *cedar.Policy]) int {
	count := 0
	for range it {
		count++
	}
	return count
}
func TestGetPolicySet(t *testing.T) {
	ctx := context.Background()
	db := setupDB(t, true)
	defer db.Close()
	ds := NewSqlDatastore(db)
	ps, err := ds.GetPolicySet(ctx, nil)
	must.NoError(t, err)
	must.Eq(t, 1, count(ps.All()))
	p := ps.Get("policy0")
	data := p.MarshalCedar()
	expectedPolicy := `permit (
    principal == User::"alice",
    action == Action::"view",
    resource in Album::"jane_vacation"
);`
	must.Eq(t, expectedPolicy, string(data))
	ps2, err := ds.GetPolicySet(ctx, &fencev1.UID{Type: "Tenant", Id: "1"})
	must.Eq(t, 2, count(ps2.All()))
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
	em, err := ds.GetEntityMap(ctx, &fencev1.UID{Id: "1"})
	must.NoError(t, err)
	must.Eq(t, expectedMap, em)
}

func TestAddTenant(t *testing.T) {
	ctx := context.Background()
	db := setupDB(t, true)
	defer db.Close()
	ds := NewSqlDatastore(db)
	fenceTenant := &fencev1.Tenant{
		Id:   &fencev1.UID{Type: "Tenant", Id: "2"},
		Name: "Test",
	}
	err := ds.AddTenant(ctx, fenceTenant)
	must.NoError(t, err)
	var tenant models.Tenant
	err = db.NewSelect().Model(&tenant).Where("id = ?", "2").Scan(ctx)
	must.NoError(t, err)
	must.Eq(t, "2", tenant.ID)
	must.Eq(t, "Test", tenant.Name)
}
func TestListTenants(t *testing.T) {
	ctx := context.Background()
	db := setupDB(t, true)
	defer db.Close()
	ds := NewSqlDatastore(db)
	records, err := ds.ListTenants(ctx)
	must.NoError(t, err)
	expected := []*fencev1.Tenant{
		{Id: &fencev1.UID{Type: "Tenant", Id: "1"}, Name: "FixtureTenant"},
	}
	must.Eq(t, expected, records, must.Cmp(protocmp.Transform()))
}
