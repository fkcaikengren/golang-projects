package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"go-oj/internal/handler"
	"go-oj/internal/model"
	"go-oj/internal/service"
)

type stubAdminAuthService struct {
	loginFn        func(ctx context.Context, input service.AdminLoginInput) (*service.AdminAuthResult, error)
	authenticateFn func(ctx context.Context, token string) (*model.AdminUser, error)
}

func (s stubAdminAuthService) Login(ctx context.Context, input service.AdminLoginInput) (*service.AdminAuthResult, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, input)
	}
	return &service.AdminAuthResult{
		Token:     "admin-token",
		AdminUser: &model.AdminUser{BaseModel: model.BaseModel{ID: 1}, Email: input.Email, DisplayName: "admin"},
	}, nil
}

func (s stubAdminAuthService) Authenticate(ctx context.Context, token string) (*model.AdminUser, error) {
	if s.authenticateFn != nil {
		return s.authenticateFn(ctx, token)
	}
	return nil, service.ErrInvalidToken
}

func newAdminTestRouter(t *testing.T, authz *service.AdminAuthorizer) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)
	return New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
		handler.NewAdminAuthHandler(stubAdminAuthService{
			authenticateFn: func(ctx context.Context, token string) (*model.AdminUser, error) {
				switch token {
				case "Bearer assistant-token":
					return &model.AdminUser{
						BaseModel:   model.BaseModel{ID: 2},
						Email:       "assistant@go-oj.dev",
						DisplayName: "Assistant",
						Status:      service.AdminUserStatusActive,
					}, nil
				case "Bearer admin-token":
					return &model.AdminUser{
						BaseModel:   model.BaseModel{ID: 1},
						Email:       "admin@go-oj.dev",
						DisplayName: "Super Admin",
						Status:      service.AdminUserStatusActive,
					}, nil
				default:
					return nil, service.ErrInvalidToken
				}
			},
		}),
		stubAdminAuthService{
			authenticateFn: func(ctx context.Context, token string) (*model.AdminUser, error) {
				switch token {
				case "Bearer assistant-token":
					return &model.AdminUser{
						BaseModel:   model.BaseModel{ID: 2},
						Email:       "assistant@go-oj.dev",
						DisplayName: "Assistant",
						Status:      service.AdminUserStatusActive,
					}, nil
				case "Bearer admin-token":
					return &model.AdminUser{
						BaseModel:   model.BaseModel{ID: 1},
						Email:       "admin@go-oj.dev",
						DisplayName: "Super Admin",
						Status:      service.AdminUserStatusActive,
					}, nil
				default:
					return nil, service.ErrInvalidToken
				}
			},
		},
		handler.NewAdminHandler(),
		authz,
	)
}

func TestAdminRoutesRequireAuthorization(t *testing.T) {
	t.Parallel()

	authz, err := service.NewInMemoryAdminAuthorizer()
	if err != nil {
		t.Fatalf("NewInMemoryAdminAuthorizer() error = %v", err)
	}
	r := newAdminTestRouter(t, authz)

	req := httptest.NewRequest(http.MethodGet, "/admin/settings", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("GET /admin/settings status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAdminRoutesEnforcePermissions(t *testing.T) {
	t.Parallel()

	authz, err := service.NewInMemoryAdminAuthorizer()
	if err != nil {
		t.Fatalf("NewInMemoryAdminAuthorizer() error = %v", err)
	}
	if err := authz.AssignRole("1", service.AdminRoleAdmin); err != nil {
		t.Fatalf("AssignRole(admin) error = %v", err)
	}
	if err := authz.AssignRole("2", service.AdminRoleAssistant); err != nil {
		t.Fatalf("AssignRole(assistant) error = %v", err)
	}

	r := newAdminTestRouter(t, authz)

	t.Run("assistant forbidden for settings write", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/admin/settings", nil)
		req.Header.Set("Authorization", "Bearer assistant-token")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("PATCH /admin/settings status = %d, want %d", rec.Code, http.StatusForbidden)
		}
	})

	t.Run("assistant allowed for problems read", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/problems", nil)
		req.Header.Set("Authorization", "Bearer assistant-token")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /admin/problems status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("admin allowed for settings write", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/admin/settings", nil)
		req.Header.Set("Authorization", "Bearer admin-token")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("PATCH /admin/settings status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

func TestAdminAuthRouteIsRegistered(t *testing.T) {
	t.Parallel()

	authz, err := service.NewInMemoryAdminAuthorizer()
	if err != nil {
		t.Fatalf("NewInMemoryAdminAuthorizer() error = %v", err)
	}

	r := newAdminTestRouter(t, authz)
	req := httptest.NewRequest(http.MethodPost, "/admin/login", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code == http.StatusNotFound {
		t.Fatalf("POST /admin/login returned 404, route not registered")
	}
}

func TestAdminDashboardReturnsAdminContext(t *testing.T) {
	t.Parallel()

	authz, err := service.NewInMemoryAdminAuthorizer()
	if err != nil {
		t.Fatalf("NewInMemoryAdminAuthorizer() error = %v", err)
	}
	if err := authz.AssignRole("1", service.AdminRoleAdmin); err != nil {
		t.Fatalf("AssignRole(admin) error = %v", err)
	}

	r := newAdminTestRouter(t, authz)
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer admin-token")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /admin status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), "\"email\":\"admin@go-oj.dev\"") {
		t.Fatalf("dashboard response missing admin email: %s", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "\"display_name\":\"Super Admin\"") {
		t.Fatalf("dashboard response missing admin display name: %s", rec.Body.String())
	}
}
