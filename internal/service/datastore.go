package service

import (
	"context"

	"github.com/cedar-policy/cedar-go"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type DataStore interface {
	addPolicy(ctx context.Context, id string, content string) error
	deletePolicy(ctx context.Context, id string) error
	getPolicySet(ctx context.Context) (*cedar.PolicySet, error)
	getPolicies(ctx context.Context) ([]*fencev1.Policy, error)
	getPolicy(ctx context.Context, id string) (*cedar.Policy, error)
	getEntities(ctx context.Context) ([]*fencev1.Entity, error)
	getEntityMap(ctx context.Context) (cedar.EntityMap, error)
	addEntity(ctx context.Context, entity *fencev1.Entity) error
	deleteEntity(ctx context.Context, uid *fencev1.UID) error
}
