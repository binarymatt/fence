package providers

import (
	"errors"
	"fmt"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

var (
	ErrInvalidPrincipal = errors.New("invalid principal")
)

func NewAuthzError(principal, action, resource *fencev1.UID) FenceAuthzError {
	message := fmt.Sprintf("%s not allowed to %s on %s", UIDToString(principal), UIDToString(action), UIDToString(resource))
	return FenceAuthzError{message: message}
}

type FenceAuthzError struct {
	message string
}

func (az FenceAuthzError) Error() string {
	return az.message
}
