package database

import (
	"context"
	"fmt"

	"github.com/cedar-policy/cedar-go"
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/internal/translation"
)

const (
	policyPrefix = "policy"
	entityPrefix = "entity"
)

var _ DataStore = (*badgerDataStore)(nil)

func entityId(typ, id string) []byte {
	return fmt.Appendf(nil, "%s_%s_%s", entityPrefix, typ, id)
}
func policyId(id string) []byte {
	return fmt.Appendf(nil, "%s_%s", policyPrefix, id)
}
func NewBaderStore(db *badger.DB) *badgerDataStore {
	return &badgerDataStore{db}
}

type badgerDataStore struct {
	db *badger.DB
}

func (bds *badgerDataStore) AddPolicy(ctx context.Context, id string, content string) error {
	p := fencev1.Policy{
		Id:         id,
		Definition: content,
	}
	data, err := proto.Marshal(&p)
	if err != nil {
		return err
	}
	return bds.db.Update(func(tx *badger.Txn) error {
		err := tx.Set(policyId(id), data)
		return err
	})
}

func (bds *badgerDataStore) DeletePolicy(ctx context.Context, id string) error {
	return bds.db.Update(func(tx *badger.Txn) error {
		return tx.Delete(policyId(id))
	})
}

func (bds *badgerDataStore) GetPolicySet(ctx context.Context, tenant *fencev1.UID) (*cedar.PolicySet, error) {
	ps := cedar.NewPolicySet()
	policies, err := bds.GetPolicies(ctx, tenant)
	for _, p := range policies {
		var policy cedar.Policy
		if err := policy.UnmarshalCedar([]byte(p.Definition)); err != nil {
			return nil, err
		}
		ps.Add(cedar.PolicyID(p.Id), &policy)
	}
	return ps, err
}

func (bds *badgerDataStore) GetPolicies(ctx context.Context, tenant *fencev1.UID) ([]*fencev1.Policy, error) {
	policies := []*fencev1.Policy{}
	err := bds.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(policyPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			// k := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			var p fencev1.Policy
			if err := proto.Unmarshal(val, &p); err != nil {
				return err
			}
			policies = append(policies, &p)
		}
		return nil
	})
	return policies, err
}
func (bds *badgerDataStore) GetPolicy(ctx context.Context, id string) (*cedar.Policy, error) {
	var policy cedar.Policy
	err := bds.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		data, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		var p fencev1.Policy
		if err := proto.Unmarshal(data, &p); err != nil {
			return err
		}
		if err := policy.UnmarshalCedar([]byte(p.Definition)); err != nil {
			return err
		}
		return nil
	})
	return &policy, err
}
func (bds *badgerDataStore) GetEntity(ctx context.Context, tenant, uid *fencev1.UID) (*fencev1.Entity, error) {
	var entity fencev1.Entity
	err := bds.db.View(func(txn *badger.Txn) error {

		id := entityId(uid.GetType(), uid.GetId())
		item, err := txn.Get([]byte(id))
		if err != nil {
			return nil
		}
		data, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		if err := proto.Unmarshal(data, &entity); err != nil {
			return err
		}
		return nil
	})
	return &entity, err
}
func (bds *badgerDataStore) GetEntities(ctx context.Context, tenant *fencev1.UID) ([]*fencev1.Entity, error) {
	entities := []*fencev1.Entity{}
	err := bds.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := []byte(entityPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			// k := item.Key()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			var e fencev1.Entity
			if err := proto.Unmarshal(val, &e); err != nil {
				return err
			}
			entities = append(entities, &e)
		}
		return nil
	})
	return entities, err
}
func (bds *badgerDataStore) GetEntityMap(ctx context.Context, tenant *fencev1.UID) (cedar.EntityMap, error) {
	entities, err := bds.GetEntities(ctx, tenant)
	if err != nil {
		return nil, err
	}
	em := cedar.EntityMap{}
	for _, e := range entities {
		id := translation.FenceToCedarUID(e.GetUid())
		parents := make([]cedar.EntityUID, len(e.GetParents()))

		for i, p := range e.GetParents() {
			parents[i] = translation.FenceToCedarUID(p)
		}
		entity := cedar.Entity{
			UID:        id,
			Parents:    cedar.NewEntityUIDSet(parents...),
			Tags:       translation.MapToRecord(e.GetTags()),
			Attributes: translation.MapToRecord(e.GetAttributes()),
		}
		em[entity.UID] = entity
	}
	return em, nil
}
func (bds *badgerDataStore) AddEntities(ctx context.Context, tenant *fencev1.UID, entities ...*fencev1.Entity) error {
	for _, entity := range entities {
		data, err := proto.Marshal(entity)
		if err != nil {
			return err
		}
		id := entityId(entity.Uid.Type, entity.Uid.Id)
		err = bds.db.Update(func(tx *badger.Txn) error {

			return tx.Set(id, data)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (bds *badgerDataStore) DeleteEntity(ctx context.Context, tenant, uid *fencev1.UID) error {
	id := entityId(uid.Type, uid.Id)
	return bds.db.Update(func(tx *badger.Txn) error {
		return tx.Delete(id)
	})
}

func (bds *badgerDataStore) AddTenant(ctx context.Context, tenant *fencev1.Tenant) error {
	return nil
}

func (bds *badgerDataStore) ListTenants(ctx context.Context) ([]*fencev1.Tenant, error) {
	return nil, nil
}
func (bds *badgerDataStore) GetTenant(ctx context.Context, uid *fencev1.UID) (*fencev1.Tenant, error) {
	return nil, nil
}
