package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"log/slog"
)

// Redirect выполняет редирект на сохранённый URL по alias.
func Redirect(log *slog.Logger, s storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			http.Error(w, "alias is required", http.StatusBadRequest)
			return
		}

		url, err := s.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Error("get url", sl.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
