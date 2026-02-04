package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RealIP возвращает middleware, которое выставляет реальный IP клиента из X-Forwarded-For / X-Real-IP.
func RealIP(next http.Handler) http.Handler {
	return middleware.RealIP(next)
}
