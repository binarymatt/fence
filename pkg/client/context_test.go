package client

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
)

func TestClientFromContext(t *testing.T) {
	c := &Client{}
	ctx := context.WithValue(context.Background(), clientContextKey, c)
	cl := ClientFromContext(ctx)
	must.Eq(t, c, cl)

	cl = ClientFromContext(context.Background())
	must.Nil(t, cl)
}
func TestContexWithClient(t *testing.T) {
	ctx := ContextWithClient(context.Background(), &Client{})
	val := ctx.Value(clientContextKey)
	cl, ok := val.(*Client)
	must.True(t, ok)
	must.Eq(t, &Client{}, cl)
}
