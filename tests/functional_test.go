//go:build functional

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"url-shortener/internal/httpserver"
	"url-shortener/internal/storage/sqlite"
)

func TestFunctional_SaveRedirectDelete(t *testing.T) {
	log := slog.Default()
	store, err := sqlite.New(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	srv := httpserver.New(httpserver.Config{
		Addr:        ":0",
		Token:       "test-token",
	}, log, store)
	// Use httptest with router directly to avoid binding to port
	r := httpserver.NewRouter(httpserver.RouterDeps{Log: log, Store: store, Token: "test-token"})

	// Save
	body := map[string]string{"url": "https://example.com", "alias": "ex"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/url", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("POST /url status = %d; want %d", rr.Code, http.StatusCreated)
	}

	// Redirect (router parses /url/ex -> alias=ex)
	req = httptest.NewRequest(http.MethodGet, "/url/ex", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusFound {
		t.Errorf("GET /url/ex status = %d; want %d", rr.Code, http.StatusFound)
	}
	if loc := rr.Header().Get("Location"); loc != "https://example.com" {
		t.Errorf("Location = %q; want https://example.com", loc)
	}

	// Delete (with auth)
	req = httptest.NewRequest(http.MethodDelete, "/url/ex", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Errorf("DELETE /url/ex status = %d; want %d", rr.Code, http.StatusNoContent)
	}

	// Redirect after delete -> 404
	req = httptest.NewRequest(http.MethodGet, "/url/ex", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("GET /url/ex after delete status = %d; want %d", rr.Code, http.StatusNotFound)
	}

	_ = srv
}
