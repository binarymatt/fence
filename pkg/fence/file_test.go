package fence

import (
	"context"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"
	"github.com/spf13/afero"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

func createTestFS(t *testing.T) afero.Fs {
	t.Helper()
	fs := afero.NewMemMapFs()
	policies := `permit (
	principal == User::"bob",
	action == Action::"view",
	resource in Album::"vacation"
);
permit (
	principal == User::"jane",
	action == Action::"create",
	resource in Album::"home"
);
`

	entitiesJSON := `[
  {
    "uid": { "type": "User", "id": "bob" },
    "attrs": { "age": 18 },
  	"parents": []
  },
  {
    "uid": { "type": "User", "id": "jane" },
    "attrs": { "age": 40 },
  	"parents": []
  },
  {
    "uid": { "type": "Photo", "id": "VacationPhoto94.jpg" },
    "attrs": {},
    "parents": [{ "type": "Album", "id": "vacation" }]
  }
]`
	err := afero.WriteFile(fs, "./policies.cedar", []byte(policies), 0644)
	must.NoError(t, err)
	err = afero.WriteFile(fs, "./entities.json", []byte(entitiesJSON), 0644)
	must.NoError(t, err)

	return fs
}
func TestNewFileState(t *testing.T) {
	fs := createTestFS(t)
	bob := cedar.Entity{
		UID:        cedar.NewEntityUID(cedar.EntityType("User"), cedar.String("bob")),
		Attributes: cedar.NewRecord(cedar.RecordMap{"age": cedar.Long(18)}),
	}
	jane := cedar.Entity{
		UID:        cedar.NewEntityUID(cedar.EntityType("User"), cedar.String("jane")),
		Attributes: cedar.NewRecord(cedar.RecordMap{"age": cedar.Long(40)}),
	}
	photo := cedar.Entity{
		UID:        cedar.NewEntityUID(cedar.EntityType("Photo"), cedar.String("VacationPhoto94.jpg")),
		Attributes: cedar.NewRecord(cedar.RecordMap{}),
		Parents:    cedar.NewEntityUIDSet(cedar.NewEntityUID(cedar.EntityType("Album"), cedar.String("vacation"))),
	}
	state, err := NewFileState(fs, "./policies.cedar", "./entities.json")
	must.NoError(t, err)
	expectedEntities := cedar.EntityMap{
		bob.UID:   bob,
		jane.UID:  jane,
		photo.UID: photo,
	}
	must.Eq(t, expectedEntities, state.entities)
	must.MapLen(t, 2, state.ps.Map())
}
func TestFileIsAllowed(t *testing.T) {
	fs := createTestFS(t)
	state, err := NewFileState(fs, "policies.cedar", "entities.json")
	must.NoError(t, err)
	bob := &fencev1.UID{
		Type: "User",
		Id:   "bob",
	}
	jane := &fencev1.UID{
		Type: "User",
		Id:   "jane",
	}
	resource := &fencev1.UID{
		Type: "Photo",
		Id:   "VacationPhoto94.jpg",
	}
	action := &fencev1.UID{
		Type: "Action",
		Id:   "view",
	}
	err = state.IsAllowed(context.Background(), bob, action, resource)
	must.NoError(t, err)

	err = state.IsAllowed(context.Background(), jane, action, resource)
	must.Error(t, err)
	var fe FenceAuthzError
	must.ErrorAs(t, err, &fe)
	must.Eq(t, `User::"jane" not allowed to Action::"view" on Photo::"VacationPhoto94.jpg"`, fe.Error())
}
