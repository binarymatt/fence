package fence

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"

	fencev1 "github.com/binarymatt/fence/gen/fence/v1"
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
func TestPrincipalFromContext(t *testing.T) {
	p := &fencev1.UID{Type: "user", Id: "bob"}
	ctx := context.WithValue(context.Background(), principaContextKey, p)
	pl, err := PrincipalFromContext(ctx)
	must.NoError(t, err)
	must.Eq(t, p, pl)

	pl, err = PrincipalFromContext(context.Background())
	must.ErrorIs(t, err, ErrInvalidPrincipal)
}
func TestContexWithPrincipal(t *testing.T) {
	p := &fencev1.UID{Type: "user", Id: "bob"}
	ctx := ContextWithPrincipal(context.Background(), p)
	val := ctx.Value(principaContextKey)
	cl, ok := val.(*fencev1.UID)
	must.True(t, ok)
	must.Eq(t, p, cl)
}
