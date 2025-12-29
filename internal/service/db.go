package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cedar-policy/cedar-go"
	"github.com/uptrace/bun"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrPolicyNotFound      = errors.New("policy not found")
	ErrPolicyAlreadyExists = errors.New("policy already exists")
)

var _ DataStore = (*sqlDataStore)(nil)

func NewSqlDatastore(db *bun.DB) *sqlDataStore {
	return &sqlDataStore{db}
}

type sqlDataStore struct {
	db *bun.DB
}

func (s *sqlDataStore) addPolicy(ctx context.Context, id string, content string) error {
	p := &Policy{
		ID:      id,
		Content: content,
	}
	_, err := s.db.NewInsert().Model(p).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed: policies.id") {
			return fmt.Errorf("%w: policy %s", ErrPolicyAlreadyExists, id)
		}
		return err
	}
	return nil
}

func (s *sqlDataStore) deletePolicy(ctx context.Context, id string) error {
	res, err := s.db.NewDelete().Model(&Policy{ID: id}).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf(`%s: %w`, id, ErrPolicyNotFound)
	}
	return nil
}
func (s *sqlDataStore) getPolicySet(ctx context.Context) (*cedar.PolicySet, error) {

	policies, err := s.getPolicies(ctx)
	if err != nil {
		return nil, err
	}
	ps := cedar.NewPolicySet()
	for _, rawPolicy := range policies {
		var policy cedar.Policy
		if err := policy.UnmarshalCedar([]byte(rawPolicy.Definition)); err != nil {
			return nil, err
		}
		ps.Add(cedar.PolicyID(rawPolicy.Id), &policy)
	}
	return ps, nil
}
func (s *sqlDataStore) getPolicies(ctx context.Context) ([]*fencev1.Policy, error) {
	var policies []Policy
	if err := s.db.NewSelect().Model(&policies).Scan(ctx); err != nil {
		return nil, err
	}
	pols := make([]*fencev1.Policy, len(policies))
	for i, policy := range policies {
		pols[i] = policy.ToProto()
	}
	return pols, nil
}
func (s *sqlDataStore) getPolicy(ctx context.Context, id string) (*cedar.Policy, error) {
	var policy Policy
	if err := s.db.NewSelect().Model(&policy).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	var cPolicy cedar.Policy
	if err := cPolicy.UnmarshalCedar([]byte(policy.Content)); err != nil {
		return nil, err
	}
	return &cPolicy, nil
}

func (s *sqlDataStore) getEntities(ctx context.Context) ([]*fencev1.Entity, error) {
	var entities []Entity
	if err := s.db.NewSelect().Model(&entities).Scan(ctx); err != nil {
		return nil, err
	}
	// ents := make([]cedar.Entity, len(entities))
	ents := make([]*fencev1.Entity, len(entities))
	for i, e := range entities {
		ents[i] = e.ToProto()
	}
	return ents, nil
}
func (s *sqlDataStore) getEntityMap(ctx context.Context) (cedar.EntityMap, error) {

	entities, err := s.getEntities(ctx)
	if err != nil {
		return nil, err
	}
	em := cedar.EntityMap{}
	for _, e := range entities {
		id := cedar.NewEntityUID(cedar.EntityType(e.Uid.Type), cedar.String(e.Uid.Id))
		parentRecords := []cedar.EntityUID{}
		for _, uid := range e.Parents {
			parentRecords = append(parentRecords, cedar.NewEntityUID(cedar.EntityType(uid.Type), cedar.String(uid.Id)))
		}
		parents := cedar.NewEntityUIDSet(parentRecords...)
		ent := cedar.Entity{
			UID:     id,
			Parents: parents,
			//Attributes: e.Attributes,
			//Tags:       e.Tags,
		}
		em[id] = ent
	}
	return em, nil
}

func (s *sqlDataStore) addEntity(ctx context.Context, e *fencev1.Entity) error {
	uids := make([]UID, len(e.Parents))
	for _, uid := range e.Parents {

		uids = append(uids, UID{ID: string(uid.Id), Type: string(uid.Type)})
	}
	entity := &Entity{
		ID:      string(e.Uid.Id),
		Type:    string(e.Uid.Type),
		Parents: uids,
		//Attributes: e.Attributes,
		//Tags:       e.Tags,
	}
	_, err := s.db.NewInsert().Model(entity).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed: entities.id, entities.type") {
			return fmt.Errorf(`%w: %s::"%s"`, ErrEntityAlreadyExists, entity.Type, entity.ID)
		}
		return err
	}
	return nil
}
func (s *sqlDataStore) deleteEntity(ctx context.Context, e *fencev1.UID) error {
	id := string(e.Id)
	typ := string(e.Type)
	res, err := s.db.NewDelete().Model(&Entity{ID: id, Type: typ}).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return fmt.Errorf(`%s::"%s": %w`, typ, id, ErrEntityNotFound)
	}
	return nil
}
