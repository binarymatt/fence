package agent

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"golang.org/x/sync/errgroup"

	"github.com/binarymatt/fence/gen/fence/v1/fencev1connect"
	"github.com/binarymatt/fence/internal/service"
)

func initDB(ctx context.Context, db *bun.DB) error {

	_, err := db.NewCreateTable().Model((*service.Entity)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateTable().Model((*service.Policy)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
func New(ctx context.Context, cfg *Config) (*agent, error) {
	sqlDB, err := sql.Open(sqliteshim.ShimName, cfg.DBPath)
	if err != nil {
		return nil, err
	}
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	if err := initDB(ctx, db); err != nil {
		return nil, err
	}
	svc := service.New(db)
	return &agent{cfg, svc}, nil
}

type agent struct {
	cfg     *Config
	service *service.Service
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("processing request", "method", r.Method, "uri", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, connect-protocol-version")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func (a *agent) Run(ctx context.Context) error {
	slog.Info("starting agent run", "address", a.cfg.ListenAddress)
	eg, ctx := errgroup.WithContext(ctx)
	mux := http.NewServeMux()
	mux.Handle(fencev1connect.NewFenceServiceHandler(a.service))
	mux.Handle(fencev1connect.NewFenceAdminServiceHandler(a.service))
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		slog.Info("health check")
		w.Write([]byte("ok"))
	})
	p := new(http.Protocols)
	p.SetHTTP1(true)
	// For gRPC clients, it's convenient to support HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      a.cfg.ListenAddress,
		Handler:   loggingMiddleware(corsMiddleware(mux)),
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
