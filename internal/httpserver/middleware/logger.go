package middleware

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger логирует метод, путь, IP, длительность и статус каждого запроса.
func Logger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID := middleware.GetReqID(r.Context())
			path := r.URL.Path
			clientIP := r.RemoteAddr
			method := r.Method

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			log.Info("request",
				slog.String("request_id", reqID),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", clientIP),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
