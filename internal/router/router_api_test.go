package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-oj/internal/handler"
	"go-oj/internal/service"
)

func newTestAPIRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	adminAuthz, err := service.NewInMemoryAdminAuthorizer()
	if err != nil {
		panic(err)
	}
	return New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
		handler.NewAdminAuthHandler(stubAdminAuthService{}),
		stubAdminAuthService{},
		handler.NewAdminHandler(),
		adminAuthz,
	)
}

func TestProblemSetRoutesAreRegistered(t *testing.T) {
	r := newTestAPIRouter()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "list", method: http.MethodGet, path: "/api/v1/problem-sets"},
		{name: "detail", method: http.MethodGet, path: "/api/v1/problem-sets/intro"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code == http.StatusNotFound {
				t.Fatalf("%s %s returned 404, route not registered", tt.method, tt.path)
			}
			if got := rec.Header().Get("Content-Type"); got == "" {
				t.Fatalf("expected JSON response content type, got empty")
			}
		})
	}
}

func TestProblemRoutesAreRegistered(t *testing.T) {
	r := newTestAPIRouter()

	tests := []struct {
		name string
		path string
	}{
		{name: "list", path: "/api/v1/problems"},
		{name: "detail", path: "/api/v1/problems/two-sum"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code == http.StatusNotFound {
				t.Fatalf("GET %s returned 404, route not registered", tt.path)
			}
			if got := rec.Header().Get("Content-Type"); got == "" {
				t.Fatalf("expected JSON response content type, got empty")
			}
		})
	}
}

func TestSubmissionRoutesAreRegistered(t *testing.T) {
	r := newTestAPIRouter()

	t.Run("submit", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodPost,
			"/api/v1/submissions",
			bytes.NewBufferString(`{"user_id":1,"problem_id":2,"language":"go","code":"package main"}`),
		)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code == http.StatusNotFound {
			t.Fatalf("POST /api/v1/submissions returned 404, route not registered")
		}
	})

	t.Run("list my submissions validates query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/submissions", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("GET /api/v1/submissions status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("list problem submissions validates query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/problems/two-sum/submissions", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("GET /api/v1/problems/:slug/submissions status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
