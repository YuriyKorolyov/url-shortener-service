package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"log/slog"
)

// Delete удаляет URL по alias.
func Delete(log *slog.Logger, s storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			http.Error(w, "alias is required", http.StatusBadRequest)
			return
		}

		err := s.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			log.Error("delete url", sl.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
