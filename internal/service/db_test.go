package service

import (
	"context"
	"testing"

	"github.com/cedar-policy/cedar-go"
	"github.com/shoenig/test/must"
)

func TestGetPolicySet(t *testing.T) {
	s, db := setupTest(t, true)
	defer db.Close()
	ps, err := s.getPolicySet(context.Background())
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
	s, db := setupTest(t, true)
	defer db.Close()

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
	em, err := s.getEntityMap(ctx)
	must.NoError(t, err)
	must.Eq(t, expectedMap, em)
}
