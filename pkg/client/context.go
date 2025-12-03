package client

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type contextKeyType struct {
	name string
}

var clientContextKey = contextKeyType{"client"}
var principaContextKey = contextKeyType{"principal"}

func ClientFromContext(ctx context.Context) *client {
	c := ctx.Value(clientContextKey)
	cl, ok := c.(*client)
	if !ok {
		return nil
	}
	return cl
}

func ContextWithClient(ctx context.Context, cl *client) context.Context {
	return context.WithValue(ctx, clientContextKey, cl)
}

func PrincipalFromContext(ctx context.Context) (*fencev1.UID, error) {
	principal, ok := ctx.Value(principaContextKey).(*fencev1.UID)
	if !ok {
		return principal, ErrInvalidPrincipal
	}
	return principal, nil
}

func ContextWithPrincipal(ctx context.Context, principal *fencev1.UID) context.Context {
	return context.WithValue(ctx, principaContextKey, principal)
}
