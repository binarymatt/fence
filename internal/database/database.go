package database

import (
	"context"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type DataStore interface {
	AddPolicy(ctx context.Context, id string, content string) error
	DeletePolicy(ctx context.Context, id string) error
	GetPolicySet(ctx context.Context, tenant *fencev1.UID) (*cedar.PolicySet, error)
	GetPolicies(ctx context.Context, tenant *fencev1.UID) ([]*fencev1.Policy, error)
	GetPolicy(ctx context.Context, id string) (*cedar.Policy, error)
	GetEntities(ctx context.Context, parent *fencev1.UID) ([]*fencev1.Entity, error)
	GetEntityMap(ctx context.Context) (cedar.EntityMap, error)
	// AddEntity(ctx context.Context, entity *fencev1.Entity) error
	AddEntities(ctx context.Context, entities ...*fencev1.Entity) error
	DeleteEntity(ctx context.Context, uid *fencev1.UID) error
	GetEntity(ctx context.Context, uid *fencev1.UID) (*fencev1.Entity, error)
	AddTenant(ctx context.Context, uid *fencev1.UID) error
	ListTenants(ctx context.Context) ([]*fencev1.UID, error)
}
