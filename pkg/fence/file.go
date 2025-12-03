package fence

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/cedar-policy/cedar-go"
	"github.com/spf13/afero"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var _ FenceState = (*FileState)(nil)

type FileState struct {
	entityPath string
	policyPath string
	ps         *cedar.PolicySet
	entities   cedar.EntityMap
	fs         afero.Fs
}

func (s *FileState) IsAllowed(ctx context.Context, principal *fencev1.UID, action *fencev1.UID, resource *fencev1.UID) error {
	cPrincipal := uidToCedar(principal)
	cResource := uidToCedar(resource)
	cAction := uidToCedar(action)
	req := cedar.Request{
		Principal: cPrincipal,
		Resource:  cResource,
		Action:    cAction,
	}
	decision, _ := cedar.Authorize(s.ps, s.entities, req)

	if !decision {
		return NewAuthzError(principal, action, resource)
	}
	return nil
}
func (s *FileState) refresh() error {
	policyData, err := afero.ReadFile(s.fs, s.policyPath)
	if err != nil {
		slog.Error("failed to read file for policies", "path", s.policyPath)
		return err
	}
	entityData, err := afero.ReadFile(s.fs, s.entityPath)
	if err != nil {
		slog.Error("failed to read file for entities", "path", s.entityPath)
		return err
	}
	ps, err := cedar.NewPolicySetFromBytes(s.policyPath, policyData)
	if err != nil {
		slog.Error("failed to read policy file", "error", err, "data", string(policyData))
		return err
	}
	s.ps = ps

	var entities cedar.EntityMap
	if err = json.Unmarshal([]byte(entityData), &entities); err != nil {
		return err
	}
	s.entities = entities
	return nil
}
func (s *FileState) Refresh(context.Context) error {
	return nil
}
func NewFileState(fs afero.Fs, policyPath, entityPath string) (*FileState, error) {
	state := &FileState{
		policyPath: policyPath,
		entityPath: entityPath,
		fs:         fs,
	}
	if err := state.refresh(); err != nil {
		return nil, err
	}
	return state, nil
}
