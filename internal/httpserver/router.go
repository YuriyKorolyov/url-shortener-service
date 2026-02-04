package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"url-shortener/internal/httpserver/handlers"
	"url-shortener/internal/httpserver/middleware"
	"url-shortener/internal/storage"

	"log/slog"
)

type RouterDeps struct {
	Log   *slog.Logger
	Store storage.URLStorage
	Token string
}

// NewRouter собирает Chi роутер с middlewares и хэндлерами.
func NewRouter(d RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger(d.Log))
	r.Use(middleware.Recover(d.Log))
	r.Use(middleware.URLFormat)

	r.Route("/url", func(r chi.Router) {
		r.Post("/", handlers.Save(d.Log, d.Store))

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(d.Token))
			r.Delete("/{alias}", handlers.Delete(d.Log, d.Store))
		})

		r.Get("/{alias}", handlers.Redirect(d.Log, d.Store))
	})

	return r
}
