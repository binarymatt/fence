package client

import (
	"context"
)

type contextKeyType struct {
	name string
}

var clientContextKey = contextKeyType{"client"}

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
