package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"url-shortener/internal/storage"
)

func TestSave(t *testing.T) {
	log := slog.Default()

	t.Run("success", func(t *testing.T) {
		mock := &storage.MockStorage{
			SaveURLFunc: func(urlToSave, alias string) (int64, error) {
				if urlToSave != "https://example.com" || alias != "ex" {
					t.Errorf("unexpected args: url=%q alias=%q", urlToSave, alias)
				}
				return 1, nil
			},
		}
		body := SaveRequest{URL: "https://example.com", Alias: "ex"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		Save(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusCreated)
		}
		var res SaveResponse
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		if res.Alias != "ex" {
			t.Errorf("alias = %q; want ex", res.Alias)
		}
	})

	t.Run("conflict alias exists", func(t *testing.T) {
		mock := &storage.MockStorage{
			SaveURLFunc: func(_, _ string) (int64, error) {
				return 0, storage.ErrURLExists
			},
		}
		body := SaveRequest{URL: "https://a.com", Alias: "a"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		Save(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusConflict {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusConflict)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mock := &storage.MockStorage{
			SaveURLFunc: func(_, _ string) (int64, error) {
				return 0, errors.New("db error")
			},
		}
		body := SaveRequest{URL: "https://a.com", Alias: "a"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		Save(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusInternalServerError)
		}
	})

	t.Run("bad request missing url", func(t *testing.T) {
		mock := &storage.MockStorage{}
		body := SaveRequest{URL: "", Alias: "x"}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		Save(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusBadRequest)
		}
	})
}
