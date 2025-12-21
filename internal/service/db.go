package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cedar-policy/cedar-go"
	"github.com/uptrace/bun"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrPolicyNotFound      = errors.New("policy not found")
	ErrPolicyAlreadyExists = errors.New("policy already exists")
)

func (s *Service) addPolicy(ctx context.Context, tx bun.Tx, id string, content string) error {
	p := &Policy{
		ID:      id,
		Content: content,
	}
	_, err := tx.NewInsert().Model(p).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed: policies.id") {
			return fmt.Errorf("%w: policy %s", ErrPolicyAlreadyExists, id)
		}
		return err
	}
	return nil
}

func (s *Service) deletePolicy(ctx context.Context, id string) error {
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
func (s *Service) getPolicySet(ctx context.Context) (*cedar.PolicySet, error) {

	policies, err := s.getPolicies(ctx)
	if err != nil {
		return nil, err
	}
	ps := cedar.NewPolicySet()
	for _, rawPolicy := range policies {
		var policy cedar.Policy
		if err := policy.UnmarshalCedar([]byte(rawPolicy.Content)); err != nil {
			return nil, err
		}
		ps.Add(cedar.PolicyID(rawPolicy.ID), &policy)
	}
	return ps, nil
}
func (s *Service) getPolicies(ctx context.Context) ([]Policy, error) {
	var policies []Policy
	if err := s.db.NewSelect().Model(&policies).Scan(ctx); err != nil {
		return nil, err
	}
	return policies, nil
}
func (s *Service) getPolicy(ctx context.Context, id string) (*Policy, error) {
	var policy Policy
	if err := s.db.NewSelect().Model(&policy).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *Service) getEntities(ctx context.Context) ([]Entity, error) {
	var entities []Entity
	if err := s.db.NewSelect().Model(&entities).Scan(ctx); err != nil {
		return nil, err
	}
	return entities, nil
}
func (s *Service) getEntityMap(ctx context.Context) (cedar.EntityMap, error) {

	var entities []Entity
	entities, err := s.getEntities(ctx)
	if err != nil {
		return nil, err
	}
	em := cedar.EntityMap{}
	for _, e := range entities {
		id := cedar.NewEntityUID(cedar.EntityType(e.Type), cedar.String(e.ID))
		parentRecords := []cedar.EntityUID{}
		for _, uid := range e.Parents {
			parentRecords = append(parentRecords, cedar.NewEntityUID(cedar.EntityType(uid.Type), cedar.String(uid.ID)))
		}
		parents := cedar.NewEntityUIDSet(parentRecords...)
		ent := cedar.Entity{
			UID:        id,
			Parents:    parents,
			Attributes: e.Attributes,
			Tags:       e.Tags,
		}
		em[id] = ent
	}
	return em, nil
}

func (s *Service) addEntity(ctx context.Context, tx bun.Tx, entity *Entity) error {
	_, err := tx.NewInsert().Model(entity).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed: entities.id, entities.type") {
			return fmt.Errorf(`%w: %s::"%s"`, ErrEntityAlreadyExists, entity.Type, entity.ID)
		}
		return err
	}
	return nil
}
func (s *Service) deleteEntity(ctx context.Context, typ, id string) error {
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
