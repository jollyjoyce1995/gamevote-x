package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gamevote-api-go/internal/server"
	"gamevote-api-go/internal/steam"
	"gamevote-api-go/internal/store"
)

func main() {
	initLogger()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.Open(ctx, env("DATABASE_URL", "postgres://gamevote:gamevote@localhost:5432/gamevote"))
	if err != nil {
		slog.Error("database setup failed", "error", err)
		os.Exit(1)
	}
	defer st.Close()

	staticDir := env("STATIC_DIR", "web")
	if _, err := os.Stat(staticDir); err != nil {
		slog.Info("not serving a frontend (STATIC_DIR not found)", "dir", staticDir)
		staticDir = ""
	}

	broker := server.NewBroker()
	srv := &http.Server{
		Addr:              ":" + env("PORT", "8080"),
		Handler:           server.New(st, broker, steam.New()).Handler(staticDir),
		ReadHeaderTimeout: 10 * time.Second,
		// No WriteTimeout: it would cut the long-lived SSE streams.
	}

	go func() {
		<-ctx.Done()
		slog.Info("shutting down")
		broker.Close()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	slog.Info("listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initLogger() {
	var level slog.Level
	_ = level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))) // empty/unknown -> info
	opts := &slog.HandlerOptions{Level: level}
	if os.Getenv("LOG_FORMAT") == "json" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)))
		return
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opts)))
}
