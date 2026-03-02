package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Code-Hex/synchro/tz"
	"github.com/cockroachdb/errors"
	slogctx "github.com/veqryn/slog-context"
	"github.com/walnuts1018/PRFExample/server/config"
	"github.com/walnuts1018/PRFExample/server/domain/webauthn"
	"github.com/walnuts1018/PRFExample/server/infra/scylladb"
	"github.com/walnuts1018/PRFExample/server/logger"
	"github.com/walnuts1018/PRFExample/server/router"
	"github.com/walnuts1018/PRFExample/server/router/handler"
	"github.com/walnuts1018/PRFExample/server/tracer"
	"github.com/walnuts1018/PRFExample/server/usecase"
	"github.com/walnuts1018/PRFExample/server/util/clock"
	"github.com/walnuts1018/PRFExample/server/util/random"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load config",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger := logger.CreateLogger(cfg.LogLevel, cfg.LogType)
	slog.SetDefault(logger)

	closeFunc, err := tracer.NewTracerProvider(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create tracer provider",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer closeFunc()

	scylla, err := scylladb.NewScyllaDB(cfg.ScyllaDB)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create ScyllaDB client",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	defer scylla.Close()

	if err := scylla.Migrate(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to migrate ScyllaDB",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	webAuthn, err := webauthn.NewWebAuthn(ctx, cfg.Server.Origin, cfg.Server.AdditionalOrigins...)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create WebAuthn service",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	usecase := usecase.NewUsecase(
		scylla,
		scylla,
		scylla,
		scylla,
		webAuthn,
		random.New(),
		clock.Default[tz.AsiaTokyo](),
	)
	handler := handler.NewHandler(usecase, random.New())
	router := router.NewRouter(cfg.Server, handler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}
	ctx = slogctx.Append(ctx, slog.String(string(semconv.ServerPortKey), fmt.Sprintf("%d", cfg.Server.Port)))

	go func() {
		slog.Info("Server is running")
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.ErrorContext(ctx, "Failed to run server", slog.String("error", err.Error()))
				os.Exit(1)
			}
		}
	}()

	<-ctx.Done()

	stop()
	slog.Info("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Failed to shutdown server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
