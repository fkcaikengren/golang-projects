package router

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-oj/internal/handler"
	"go-oj/internal/model"
	"go-oj/internal/service"
)

type stubHealthRepo struct {
	err error
}

func (s stubHealthRepo) Ping(ctx context.Context) error {
	return s.err
}

type stubAuthService struct{}

func (s stubAuthService) Register(ctx context.Context, input service.RegisterInput) (*service.AuthResult, error) {
	return &service.AuthResult{Token: "test-token", User: nil}, nil
}

func (s stubAuthService) Login(ctx context.Context, input service.LoginInput) (*service.AuthResult, error) {
	return &service.AuthResult{Token: "test-token", User: nil}, nil
}

type stubProblemSetService struct{}

func (s stubProblemSetService) List(ctx context.Context) ([]model.ProblemSet, error) {
	return []model.ProblemSet{}, nil
}

func (s stubProblemSetService) Detail(ctx context.Context, slug string) (*service.ProblemSetDetail, error) {
	return &service.ProblemSetDetail{}, nil
}

type stubProblemService struct{}

func (s stubProblemService) List(ctx context.Context, input service.ProblemListInput) ([]model.Problem, error) {
	return []model.Problem{}, nil
}

func (s stubProblemService) Detail(ctx context.Context, slug string) (*service.ProblemDetail, error) {
	return &service.ProblemDetail{}, nil
}

type stubSubmissionService struct{}

func (s stubSubmissionService) Submit(ctx context.Context, input service.SubmitCodeInput) (*model.Submission, error) {
	return &model.Submission{}, nil
}

func (s stubSubmissionService) ListUserSubmissions(ctx context.Context, userID uint) ([]model.Submission, error) {
	return []model.Submission{}, nil
}

func (s stubSubmissionService) ListProblemSubmissions(ctx context.Context, userID uint, problemSlug string) ([]model.Submission, error) {
	return []model.Submission{}, nil
}

func TestHealthRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
	)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
	)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /ready status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestNoRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
	)

	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET /missing status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestAuthRegisterRouteIsRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		bytes.NewBufferString(`{"email":"a@b.com","password":"123456","nickname":"ab"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/v1/auth/register status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAuthLoginRouteIsRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := New(
		handler.NewHealthHandler(service.NewHealthService("go-oj", stubHealthRepo{})),
		handler.NewAuthHandler(stubAuthService{}),
		handler.NewProblemSetHandler(stubProblemSetService{}),
		handler.NewProblemHandler(stubProblemService{}),
		handler.NewSubmissionHandler(stubSubmissionService{}),
	)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/login",
		bytes.NewBufferString(`{"email":"a@b.com","password":"123456"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /api/v1/auth/login status = %d, want %d", rec.Code, http.StatusOK)
	}
}
