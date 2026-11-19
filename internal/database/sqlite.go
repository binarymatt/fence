package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cedar-policy/cedar-go"
	"github.com/uptrace/bun"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
	"github.com/binarymatt/fence/internal/models"
	"github.com/binarymatt/fence/internal/translation"
)

var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrPolicyNotFound      = errors.New("policy not found")
	ErrPolicyAlreadyExists = errors.New("policy already exists")
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

var _ DataStore = (*sqlDataStore)(nil)

func NewSqlDatastore(db *bun.DB) *sqlDataStore {
	return &sqlDataStore{db}
}

type sqlDataStore struct {
	db *bun.DB
}

func (s *sqlDataStore) AddPolicy(ctx context.Context, id string, content string) error {
	p := &models.Policy{
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

func (s *sqlDataStore) DeletePolicy(ctx context.Context, id string) error {
	res, err := s.db.NewDelete().Model(&models.Policy{ID: id}).WherePK().Exec(ctx)
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
func (s *sqlDataStore) GetPolicySet(ctx context.Context, tenant *fencev1.UID) (*cedar.PolicySet, error) {

	policies, err := s.GetPolicies(ctx, tenant)
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
func (s *sqlDataStore) GetPolicies(ctx context.Context, tenant *fencev1.UID) ([]*fencev1.Policy, error) {
	var policies []models.Policy
	qs := s.db.NewSelect().Model(&policies).Where("tenant_id is NULL")
	if tenant != nil {
		qs = qs.WhereOr("tenant_id = ?", tenant.Id)
	}
	if err := qs.Scan(ctx); err != nil {
		return nil, err
	}
	pols := make([]*fencev1.Policy, len(policies))
	for i, policy := range policies {
		pols[i] = policy.ToProto()
	}
	return pols, nil
}
func (s *sqlDataStore) GetPolicy(ctx context.Context, id string) (*cedar.Policy, error) {
	var policy models.Policy
	if err := s.db.NewSelect().Model(&policy).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}
	var cPolicy cedar.Policy
	if err := cPolicy.UnmarshalCedar([]byte(policy.Content)); err != nil {
		return nil, err
	}
	return &cPolicy, nil
}
func (s *sqlDataStore) GetEntity(ctx context.Context, uid *fencev1.UID) (*fencev1.Entity, error) {
	if uid == nil {
		return nil, fmt.Errorf("nil entity")
	}
	var entity models.Entity
	if err := s.db.NewSelect().Model(&entity).Where("id = ? AND type =? ", uid.GetId(), uid.GetType()).Scan(ctx); err != nil {
		return nil, err
	}
	return entity.ToProto(), nil
}
func (s *sqlDataStore) GetEntities(ctx context.Context, parent *fencev1.UID) ([]*fencev1.Entity, error) {
	var entities []models.Entity
	query := s.db.NewSelect().Model(&entities)
	if parent != nil {
		// TODO: how to filter correctly
		// uid := models.Fencetodbuid(parent)

	}
	if err := query.Scan(ctx); err != nil {
		return nil, err
	}
	// ents := make([]cedar.Entity, len(entities))
	ents := make([]*fencev1.Entity, len(entities))
	for i, e := range entities {
		ents[i] = e.ToProto()
	}
	return ents, nil
}
func (s *sqlDataStore) GetEntityMap(ctx context.Context) (cedar.EntityMap, error) {

	entities, err := s.GetEntities(ctx, nil)
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

func (s *sqlDataStore) AddEntities(ctx context.Context, entities ...*fencev1.Entity) error {
	entitySlice := make([]models.Entity, len(entities))
	for i, e := range entities {
		uids := make([]models.UID, len(e.Parents))
		for _, uid := range e.Parents {
			uids = append(uids, models.UID{ID: string(uid.Id), Type: string(uid.Type)})
		}
		attrs := translation.MapToRecord(e.Attributes)
		tags := translation.MapToRecord(e.Tags)
		entitySlice[i] = models.Entity{
			ID:         string(e.Uid.Id),
			Type:       string(e.Uid.Type),
			Parents:    uids,
			Attributes: attrs,
			Tags:       tags,
		}
	}

	if _, err := s.db.NewInsert().Model(&entitySlice).Exec(ctx); err != nil {

		return err
	}
	return nil
}
func (s *sqlDataStore) DeleteEntity(ctx context.Context, e *fencev1.UID) error {
	id := string(e.Id)
	typ := string(e.Type)
	res, err := s.db.NewDelete().Model(&models.Entity{ID: id, Type: typ}).WherePK().Exec(ctx)
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

func (s *sqlDataStore) AddTenant(ctx context.Context, uid *fencev1.UID) error {
	tenant := &models.Tenant{
		ID: uid.Id,
	}
	_, err := s.db.NewInsert().Model(tenant).Exec(ctx)
	return err
}

func (s *sqlDataStore) ListTenants(ctx context.Context) ([]*fencev1.UID, error) {
	var tenants []models.Tenant
	err := s.db.NewSelect().Model(&tenants).Scan(ctx)
	if err != nil {
		return nil, err
	}
	tenantIds := make([]*fencev1.UID, len(tenants))
	for i, tenant := range tenants {
		tenantIds[i] = &fencev1.UID{
			Type: "Tenant",
			Id:   tenant.ID,
		}
	}
	return tenantIds, nil
}
