package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestID возвращает middleware, которое добавляет к запросу уникальный ID.
func RequestID(next http.Handler) http.Handler {
	return middleware.RequestID(next)
}
