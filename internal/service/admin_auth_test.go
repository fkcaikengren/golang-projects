package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

type stubAdminUserRepo struct {
	createCalls        int
	lastCreatedUser    *model.AdminUser
	getByEmailCalls    int
	updateLastLogin    int
	lastLoginUpdatedID uint
	lastLoginUpdatedTS int64

	createFn          func(ctx context.Context, user *model.AdminUser) error
	getByEmailFn      func(ctx context.Context, email string) (*model.AdminUser, error)
	getByIDFn         func(ctx context.Context, id uint) (*model.AdminUser, error)
	updateLastLoginFn func(ctx context.Context, id uint, ts int64) error
}

func (s *stubAdminUserRepo) Create(ctx context.Context, user *model.AdminUser) error {
	s.createCalls++
	s.lastCreatedUser = user
	if s.createFn != nil {
		return s.createFn(ctx, user)
	}
	return nil
}

func (s *stubAdminUserRepo) GetByEmail(ctx context.Context, email string) (*model.AdminUser, error) {
	s.getByEmailCalls++
	if s.getByEmailFn != nil {
		return s.getByEmailFn(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *stubAdminUserRepo) GetByID(ctx context.Context, id uint) (*model.AdminUser, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *stubAdminUserRepo) UpdateLastLogin(ctx context.Context, id uint, ts int64) error {
	s.updateLastLogin++
	s.lastLoginUpdatedID = id
	s.lastLoginUpdatedTS = ts
	if s.updateLastLoginFn != nil {
		return s.updateLastLoginFn(ctx, id, ts)
	}
	return nil
}

func TestAdminAuthServiceLogin(t *testing.T) {
	t.Parallel()

	const (
		secret   = "test-secret"
		tokenTTL = time.Hour
	)

	t.Run("successful login returns admin token and sanitized user", func(t *testing.T) {
		t.Parallel()

		repo := &stubAdminUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.AdminUser, error) {
				return &model.AdminUser{
					BaseModel:    model.BaseModel{ID: 9},
					Email:        email,
					PasswordHash: mustHash(t, "password123"),
					DisplayName:  "ops",
					Status:       AdminUserStatusActive,
				}, nil
			},
		}

		svc := NewAdminAuthService(repo, []byte(secret), tokenTTL)
		res, err := svc.Login(context.Background(), AdminLoginInput{
			Email:    "ops@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("Login() error = %v", err)
		}
		if res == nil || res.AdminUser == nil {
			t.Fatalf("expected admin user in result")
		}
		if res.Token == "" {
			t.Fatal("expected token")
		}
		if res.AdminUser.ID != 9 {
			t.Fatalf("expected admin user id 9, got %d", res.AdminUser.ID)
		}
		if res.AdminUser.PasswordHash != "" {
			t.Fatal("expected password hash to be scrubbed")
		}
		if repo.updateLastLogin != 1 {
			t.Fatalf("expected UpdateLastLogin called once, got %d", repo.updateLastLogin)
		}
	})

	t.Run("disabled admin cannot login", func(t *testing.T) {
		t.Parallel()

		repo := &stubAdminUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.AdminUser, error) {
				return &model.AdminUser{
					BaseModel:    model.BaseModel{ID: 7},
					Email:        email,
					PasswordHash: mustHash(t, "password123"),
					DisplayName:  "ops",
					Status:       AdminUserStatusDisabled,
				}, nil
			},
		}

		svc := NewAdminAuthService(repo, []byte(secret), tokenTTL)
		_, err := svc.Login(context.Background(), AdminLoginInput{
			Email:    "ops@example.com",
			Password: "password123",
		})
		if !errors.Is(err, ErrUserDisabled) {
			t.Fatalf("expected ErrUserDisabled, got %v", err)
		}
	})
}

func TestAdminAuthServiceAuthenticate(t *testing.T) {
	t.Parallel()

	const (
		secret   = "test-secret"
		tokenTTL = time.Hour
	)

	repo := &stubAdminUserRepo{
		getByIDFn: func(ctx context.Context, id uint) (*model.AdminUser, error) {
			return &model.AdminUser{
				BaseModel:   model.BaseModel{ID: id},
				Email:       "ops@example.com",
				DisplayName: "ops",
				Status:      AdminUserStatusActive,
			}, nil
		},
	}

	svc := NewAdminAuthService(repo, []byte(secret), tokenTTL)
	loginResult, err := svc.signToken(&model.AdminUser{
		BaseModel:   model.BaseModel{ID: 21},
		Email:       "ops@example.com",
		DisplayName: "ops",
		Status:      AdminUserStatusActive,
	})
	if err != nil {
		t.Fatalf("signToken() error = %v", err)
	}

	adminUser, err := svc.Authenticate(context.Background(), "Bearer "+loginResult)
	if err != nil {
		t.Fatalf("Authenticate() error = %v", err)
	}
	if adminUser.ID != 21 {
		t.Fatalf("expected admin user id 21, got %d", adminUser.ID)
	}
}
