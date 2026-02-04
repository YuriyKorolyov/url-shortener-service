package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"log/slog"
)

// SaveRequest — тело запроса POST /url.
type SaveRequest struct {
	URL   string `json:"url"`
	Alias string `json:"alias"`
}

// SaveResponse — ответ после сохранения.
type SaveResponse struct {
	Alias string `json:"alias"`
}

// Save сохраняет URL с заданным alias.
func Save(log *slog.Logger, s storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SaveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.URL == "" || req.Alias == "" {
			http.Error(w, "url and alias are required", http.StatusBadRequest)
			return
		}

		_, err := s.SaveURL(req.URL, req.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				http.Error(w, "alias already exists", http.StatusConflict)
				return
			}
			log.Error("save url", sl.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(SaveResponse{Alias: req.Alias})
	}
}
