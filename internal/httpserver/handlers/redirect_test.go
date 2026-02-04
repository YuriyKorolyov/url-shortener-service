package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"log/slog"

	"url-shortener/internal/storage"
)

func TestRedirect(t *testing.T) {
	log := slog.Default()

	t.Run("success", func(t *testing.T) {
		mock := &storage.MockStorage{
			GetURLFunc: func(alias string) (string, error) {
				if alias != "go" {
					return "", storage.ErrURLNotFound
				}
				return "https://go.dev", nil
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/url/go", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("alias", "go")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		Redirect(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusFound {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusFound)
		}
		if loc := rr.Header().Get("Location"); loc != "https://go.dev" {
			t.Errorf("Location = %q; want https://go.dev", loc)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &storage.MockStorage{
			GetURLFunc: func(_ string) (string, error) {
				return "", storage.ErrURLNotFound
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/url/missing", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("alias", "missing")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		Redirect(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusNotFound)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		mock := &storage.MockStorage{
			GetURLFunc: func(_ string) (string, error) {
				return "", errors.New("db error")
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/url/x", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("alias", "x")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		rr := httptest.NewRecorder()

		Redirect(log, mock).ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("status = %d; want %d", rr.Code, http.StatusInternalServerError)
		}
	})
}
