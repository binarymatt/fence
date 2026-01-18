package providers

import (
	"errors"
	"fmt"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var (
	ErrInvalidPrincipal = errors.New("invalid principal")
)

func DeniedMessage(principal, action, resource *fencev1.UID) string {
	return fmt.Sprintf("%s not allowed to %s on %s", UIDToString(principal), UIDToString(action), UIDToString(resource))
}
func NewAuthzError(principal, action, resource *fencev1.UID, internal error) FenceAuthzError {
	message := DeniedMessage(principal, action, resource)
	return FenceAuthzError{message: message}
}

type FenceAuthzError struct {
	message  string
	internal error
}

func (az FenceAuthzError) Error() string {
	return az.message
}
