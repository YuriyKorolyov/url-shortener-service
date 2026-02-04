package httpserver

import (
	"context"
	"net/http"
	"time"

	"url-shortener/internal/httpserver/handlers"
	"url-shortener/internal/httpserver/middleware"
	"url-shortener/internal/storage"

	"log/slog"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

type Config struct {
	Addr        string
	Timeout     time.Duration
	IdleTimeout time.Duration
	Token       string
}

// New создаёт HTTP-сервер с Chi роутером и middlewares.
func New(cfg Config, log *slog.Logger, s storage.URLStorage) *Server {
	r := NewRouter(RouterDeps{
		Log:   log,
		Store: s,
		Token: cfg.Token,
	})

	return &Server{
		log: log,
		httpServer: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

// Start запускает сервер (блокирующий вызов).
func (s *Server) Start() error {
	s.log.Info("http server starting", slog.String("addr", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

// Stop останавливает сервер.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
