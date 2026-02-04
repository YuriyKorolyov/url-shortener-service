package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// URLFormat — middleware для роутов с URL-параметрами.
// Парсинг выполняется chi автоматически; в хэндлерах используйте chi.URLParam(r, "alias").
func URLFormat(next http.Handler) http.Handler {
	_ = chi.URLParam // явное использование для документации
	return next
}
