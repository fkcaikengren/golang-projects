package bootstrap

import (
	"context"
	"strings"
	"testing"

	"go-oj/internal/config"
	"go-oj/internal/model"

	"gorm.io/gorm"
)

func TestDatabaseDSN(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     5432,
			Name:     "oj",
			User:     "postgres",
			Password: "Postgres",
			SSLMode:  "disable",
		},
	}

	dsn := cfg.Database.DSN()
	expectedParts := []string{
		"host=127.0.0.1",
		"port=5432",
		"user=postgres",
		"password=Postgres",
		"dbname=oj",
		"sslmode=disable",
	}

	for _, part := range expectedParts {
		if !strings.Contains(dsn, part) {
			t.Fatalf("DSN %q missing %q", dsn, part)
		}
	}
}

func TestInitDBFailure(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     1,
			Name:     "oj",
			User:     "postgres",
			Password: "Postgres",
			SSLMode:  "disable",
		},
	}

	if _, err := InitDB(cfg); err == nil {
		t.Fatal("InitDB() expected error, got nil")
	}
}

type stubBootstrapAdminUserRepo struct {
	createdUser  *model.AdminUser
	getByEmailFn func(ctx context.Context, email string) (*model.AdminUser, error)
	createFn     func(ctx context.Context, user *model.AdminUser) error
}

func (s *stubBootstrapAdminUserRepo) Create(ctx context.Context, user *model.AdminUser) error {
	s.createdUser = user
	if s.createFn != nil {
		return s.createFn(ctx, user)
	}
	return nil
}

func (s *stubBootstrapAdminUserRepo) GetByEmail(ctx context.Context, email string) (*model.AdminUser, error) {
	if s.getByEmailFn != nil {
		return s.getByEmailFn(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *stubBootstrapAdminUserRepo) GetByID(ctx context.Context, id uint) (*model.AdminUser, error) {
	return nil, gorm.ErrRecordNotFound
}

func (s *stubBootstrapAdminUserRepo) UpdateLastLogin(ctx context.Context, id uint, ts int64) error {
	return nil
}

func TestEnsureBootstrapAdminCreatesFallbackSuperAdmin(t *testing.T) {
	t.Parallel()

	repo := &stubBootstrapAdminUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*model.AdminUser, error) {
			return nil, gorm.ErrRecordNotFound
		},
		createFn: func(ctx context.Context, user *model.AdminUser) error {
			user.ID = 1
			return nil
		},
	}

	cfg := &config.Config{}
	adminUser, err := ensureBootstrapAdmin(context.Background(), repo, cfg)

	if err != nil {
		t.Fatalf("ensureBootstrapAdmin() error = %v", err)
	}
	if adminUser == nil || adminUser.ID != 1 {
		t.Fatalf("expected created bootstrap admin with id, got %#v", adminUser)
	}
	if repo.createdUser == nil {
		t.Fatal("expected bootstrap admin to be created")
	}
	if repo.createdUser.Email != "admin@go-oj.dev" {
		t.Fatalf("expected fallback email, got %q", repo.createdUser.Email)
	}
	if repo.createdUser.DisplayName != "Super Admin" {
		t.Fatalf("expected fallback display name, got %q", repo.createdUser.DisplayName)
	}
	if repo.createdUser.Status != "active" {
		t.Fatalf("expected active status, got %q", repo.createdUser.Status)
	}
	if repo.createdUser.PasswordHash == "" || repo.createdUser.PasswordHash == "Admin@123456" {
		t.Fatal("expected password hash to be generated")
	}
}

func TestEnsureBootstrapAdminSkipsExistingUser(t *testing.T) {
	t.Parallel()

	repo := &stubBootstrapAdminUserRepo{
		getByEmailFn: func(ctx context.Context, email string) (*model.AdminUser, error) {
			return &model.AdminUser{
				Email:       email,
				DisplayName: "Existing Admin",
			}, nil
		},
	}

	adminUser, err := ensureBootstrapAdmin(context.Background(), repo, &config.Config{})
	if err != nil {
		t.Fatalf("ensureBootstrapAdmin() error = %v", err)
	}
	if adminUser == nil || adminUser.Email != "admin@go-oj.dev" {
		t.Fatalf("expected existing bootstrap admin returned, got %#v", adminUser)
	}
	if repo.createdUser != nil {
		t.Fatal("expected bootstrap admin creation to be skipped")
	}
}
