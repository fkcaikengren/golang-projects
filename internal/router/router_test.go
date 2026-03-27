package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-oj/internal/handler"
	"go-oj/internal/service"
)

type stubHealthRepo struct {
	err error
}

func (s stubHealthRepo) Ping(ctx context.Context) error {
	return s.err
}

func TestHealthRoute(t *testing.T) {
	r := New(handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyRoute(t *testing.T) {
	r := New(handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})))

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /ready status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestNoRoute(t *testing.T) {
	r := New(handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})))

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET /missing status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
