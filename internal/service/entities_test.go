package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/shoenig/test/must"
	"google.golang.org/protobuf/types/known/structpb"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func TestCreateEntity(t *testing.T) {
	cases := []struct {
		name         string
		entity       *fencev1.Entity
		validate     func(t *testing.T, db *sql.DB, ent *fencev1.Entity, err error)
		loadFixutres bool
	}{
		{
			name: "basic happy path",
			entity: &fencev1.Entity{
				Uid:     &fencev1.UID{Type: "User", Id: "bob"},
				Parents: []*fencev1.UID{{Type: "Group", Id: "admins"}},
				Attributes: map[string]*structpb.Value{
					"age": structpb.NewNumberValue(40),
				},
				Tags: map[string]*structpb.Value{
					"a": structpb.NewStringValue("b"),
				},
			},
			validate: func(t *testing.T, db *sql.DB, ent *fencev1.Entity, err error) {
				var dbID, dbType, parents, tags, attrs string
				n := time.Now().UTC()
				var createdAt, updatedAt time.Time
				err = db.QueryRow("SELECT id, type, parents, attributes, tags, created_at, updated_at from entities where id =? and type = ?", "bob", "User").Scan(&dbID, &dbType, &parents, &attrs, &tags, &createdAt, &updatedAt)
				must.NoError(t, err)
				must.Eq(t, "bob", dbID)
				must.Eq(t, "User", dbType)
				must.Eq(t, `[{"id":"admins","type":"Group"}]`, parents)
				must.Eq(t, `{"age":40}`, attrs)
				must.Eq(t, `{"a":"b"}`, tags)
				must.Eq(t, createdAt.Format(time.RFC3339), n.Format(time.RFC3339))
				must.Eq(t, updatedAt.Format(time.RFC3339), n.Format(time.RFC3339))

			},
		},
		{
			name:         "existing entity",
			loadFixutres: true,
			entity: &fencev1.Entity{
				Uid:     &fencev1.UID{Type: "User", Id: "bob"},
				Parents: []*fencev1.UID{{Type: "Group", Id: "admins"}},
				Attributes: map[string]*structpb.Value{
					"age": structpb.NewNumberValue(40),
				},
				Tags: map[string]*structpb.Value{
					"a": structpb.NewStringValue("b"),
				},
			},
			validate: func(t *testing.T, db *sql.DB, ent *fencev1.Entity, err error) {
				must.Error(t, err)
				must.ErrorIs(t, err, ErrEntityAlreadyExists)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s, db := setupTest(t, tc.loadFixutres)

			req := &fencev1.CreateEntityRequest{
				Entity: tc.entity,
			}
			_, err := s.CreateEntity(context.Background(), req)
			tc.validate(t, db, tc.entity, err)

		})
	}
}

func TestDeleteEntity(t *testing.T) {
	cases := []struct {
		name     string
		id       string
		typ      string
		validate func(t *testing.T, db *sql.DB, err error)
	}{
		{
			name: "happy path",
			id:   "bob",
			typ:  "User",
			validate: func(t *testing.T, db *sql.DB, err error) {
				must.NoError(t, err)
				var count int
				db.QueryRow("select count(*) from entities").Scan(&count)
				must.Eq(t, 2, count)

			},
		},
		{
			name: "does not exist",
			id:   "jane",
			typ:  "User",
			validate: func(t *testing.T, db *sql.DB, err error) {
				must.ErrorIs(t, err, ErrEntityNotFound)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			s, db := setupTest(t, true)

			req := &fencev1.DeleteEntityRequest{
				Uid: &fencev1.UID{
					Id:   tc.id,
					Type: tc.typ,
				},
			}
			_, err := s.DeleteEntity(context.Background(), req)
			tc.validate(t, db, err)
		})
	}

}
