package agent

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
	"github.com/binarymatt/fence/internal/service"
)

type agent struct {
	cfg Config
}

func (a *agent) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	server := service.New()
	mux := http.NewServeMux()
	mux.Handle(fencev1connect.NewFenceServiceHandler(server))
	p := new(http.Protocols)
	s := http.Server{
		Addr:      a.cfg.Address,
		Handler:   mux,
		Protocols: p,
	}
	eg.Go(s.ListenAndServe)
	eg.Go(func() error {
		// Graceful shutdown
		<-ctx.Done()
		cx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()
		return s.Shutdown(cx)
	})
	return eg.Wait()
}
