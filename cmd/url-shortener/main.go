package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/config"
	"url-shortener/internal/httpserver"
	"url-shortener/internal/lib/logger/pretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	store, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error opening storage", sl.Err(err))
		os.Exit(1)
	}

	srv := httpserver.New(httpserver.Config{
		Addr:        cfg.HTTPServer.Addr,
		Timeout:     cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
		Token:       cfg.Token,
	}, log, store)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("http server", sl.Err(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.Timeout)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Error("server shutdown", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(pretty.NewHandler(os.Stdout, slog.LevelDebug))
	case envDev, envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
