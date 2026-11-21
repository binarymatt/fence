package database

import (
	"context"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type DataStore interface {
	AddPolicy(ctx context.Context, id string, content string) error
	DeletePolicy(ctx context.Context, id string) error
	GetPolicySet(ctx context.Context, tenantID *fencev1.UID) (*cedar.PolicySet, error)
	GetPolicies(ctx context.Context, tenantID *fencev1.UID) ([]*fencev1.Policy, error)
	GetPolicy(ctx context.Context, id string) (*cedar.Policy, error)
	GetEntities(ctx context.Context, tenantId *fencev1.UID) ([]*fencev1.Entity, error)
	GetEntityMap(ctx context.Context, tenant *fencev1.UID) (cedar.EntityMap, error)
	AddEntities(ctx context.Context, tenant *fencev1.UID, entities ...*fencev1.Entity) error
	DeleteEntity(ctx context.Context, tenant *fencev1.UID, uid *fencev1.UID) error
	GetEntity(ctx context.Context, tenant, uid *fencev1.UID) (*fencev1.Entity, error)
	ListTenants(ctx context.Context) ([]*fencev1.Tenant, error)
	AddTenant(ctx context.Context, tenant *fencev1.Tenant) error
	GetTenant(ctx context.Context, uid *fencev1.UID) (*fencev1.Tenant, error)
}
