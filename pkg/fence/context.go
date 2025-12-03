package fence

import (
	"context"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
)

type contextKeyType struct {
	name string
}

var clientContextKey = contextKeyType{"client"}
var principaContextKey = contextKeyType{"principal"}

func ClientFromContext(ctx context.Context) *Client {
	c := ctx.Value(clientContextKey)
	cl, ok := c.(*Client)
	if !ok {
		return nil
	}
	return cl
}

func ContextWithClient(ctx context.Context, cl *Client) context.Context {
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
