package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-oj/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type stubUserRepo struct {
	createCalls        int
	lastCreatedUser    *model.User
	getByEmailCalls    int
	updateLastLogin    int
	lastLoginUpdatedID uint
	lastLoginUpdatedTS int64

	createFn          func(ctx context.Context, user *model.User) error
	getByEmailFn      func(ctx context.Context, email string) (*model.User, error)
	updateLastLoginFn func(ctx context.Context, id uint, ts int64) error
}

func (s *stubUserRepo) Create(ctx context.Context, user *model.User) error {
	s.createCalls++
	s.lastCreatedUser = user
	if s.createFn != nil {
		return s.createFn(ctx, user)
	}
	return nil
}

func (s *stubUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	s.getByEmailCalls++
	if s.getByEmailFn != nil {
		return s.getByEmailFn(ctx, email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *stubUserRepo) UpdateLastLogin(ctx context.Context, id uint, ts int64) error {
	s.updateLastLogin++
	s.lastLoginUpdatedID = id
	s.lastLoginUpdatedTS = ts
	if s.updateLastLoginFn != nil {
		return s.updateLastLoginFn(ctx, id, ts)
	}
	return nil
}

func mustHash(t *testing.T, pw string) string {
	t.Helper()
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return string(b)
}

func parseHS256Claims(t *testing.T, token, secret string) jwt.MapClaims {
	t.Helper()

	parsed, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if tk.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("parse jwt: %v", err)
	}
	if !parsed.Valid {
		t.Fatalf("jwt invalid")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("claims type: %T", parsed.Claims)
	}
	return claims
}

func TestAuthServiceRegister(t *testing.T) {
	t.Parallel()

	const (
		secret   = "test-secret"
		tokenTTL = time.Hour
	)

	t.Run("successful registration hashes password and returns token", func(t *testing.T) {
		t.Parallel()

		repo := &stubUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			createFn: func(ctx context.Context, user *model.User) error {
				user.ID = 123
				return nil
			},
		}
		svc := NewAuthService(repo, []byte(secret), tokenTTL)

		res, err := svc.Register(context.Background(), RegisterInput{
			Email:    "a@b.com",
			Password: "password123",
			Nickname: "nick",
		})
		if err != nil {
			t.Fatalf("Register error: %v", err)
		}
		if res == nil || res.User == nil {
			t.Fatalf("expected user in result")
		}
		if repo.createCalls != 1 {
			t.Fatalf("expected 1 Create call, got %d", repo.createCalls)
		}
		if res.User.ID != 123 {
			t.Fatalf("expected ID to be set by repo, got %d", res.User.ID)
		}
		if res.User.Email != "a@b.com" {
			t.Fatalf("expected email, got %q", res.User.Email)
		}
		if res.User.Nickname != "nick" {
			t.Fatalf("expected nickname, got %q", res.User.Nickname)
		}
		if res.User.Status != "active" {
			t.Fatalf("expected status active, got %q", res.User.Status)
		}
		if res.User.LastLoginAt != 0 {
			t.Fatalf("expected last_login_at 0, got %d", res.User.LastLoginAt)
		}
		if res.User.PasswordHash != "" {
			t.Fatalf("expected returned user password hash to be scrubbed")
		}
		if repo.lastCreatedUser == nil || repo.lastCreatedUser.PasswordHash == "" {
			t.Fatalf("expected password hash to be set")
		}
		if repo.lastCreatedUser.PasswordHash == "password123" {
			t.Fatalf("expected password to be hashed, got plain text")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(repo.lastCreatedUser.PasswordHash), []byte("password123")); err != nil {
			t.Fatalf("hash does not match password: %v", err)
		}
		if res.Token == "" {
			t.Fatalf("expected token")
		}

		claims := parseHS256Claims(t, res.Token, secret)
		if claims["email"] != "a@b.com" {
			t.Fatalf("expected email claim, got %#v", claims["email"])
		}
		if claims["user_id"] == nil {
			t.Fatalf("expected user_id claim")
		}
		if claims["exp"] == nil {
			t.Fatalf("expected exp claim")
		}
	})

	t.Run("duplicate email is rejected", func(t *testing.T) {
		t.Parallel()

		repo := &stubUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{BaseModel: model.BaseModel{ID: 1}, Email: email}, nil
			},
		}
		svc := NewAuthService(repo, []byte(secret), tokenTTL)

		_, err := svc.Register(context.Background(), RegisterInput{
			Email:    "a@b.com",
			Password: "password123",
			Nickname: "nick",
		})
		if !errors.Is(err, ErrDuplicateEmail) {
			t.Fatalf("expected ErrDuplicateEmail, got %v", err)
		}
	})

	t.Run("database duplicate on create is normalized", func(t *testing.T) {
		t.Parallel()

		repo := &stubUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return nil, gorm.ErrRecordNotFound
			},
			createFn: func(ctx context.Context, user *model.User) error {
				return errors.New("duplicate key value violates unique constraint")
			},
		}
		svc := NewAuthService(repo, []byte(secret), tokenTTL)

		_, err := svc.Register(context.Background(), RegisterInput{
			Email:    "a@b.com",
			Password: "password123",
			Nickname: "nick",
		})
		if !errors.Is(err, ErrDuplicateEmail) {
			t.Fatalf("expected ErrDuplicateEmail, got %v", err)
		}
	})
}

func TestAuthServiceLogin(t *testing.T) {
	t.Parallel()

	const (
		secret   = "test-secret"
		tokenTTL = time.Hour
	)

	t.Run("successful login returns token and updates last login", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			BaseModel:    model.BaseModel{ID: 99},
			Email:        "a@b.com",
			PasswordHash: mustHash(t, "password123"),
			Nickname:     "nick",
			Status:       "active",
			LastLoginAt:  0,
		}

		repo := &stubUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return user, nil
			},
		}
		svc := NewAuthService(repo, []byte(secret), tokenTTL)

		before := time.Now().Unix()
		res, err := svc.Login(context.Background(), LoginInput{
			Email:    "a@b.com",
			Password: "password123",
		})
		after := time.Now().Unix()

		if err != nil {
			t.Fatalf("Login error: %v", err)
		}
		if res == nil || res.User == nil {
			t.Fatalf("expected user in result")
		}
		if repo.updateLastLogin != 1 {
			t.Fatalf("expected UpdateLastLogin called once, got %d", repo.updateLastLogin)
		}
		if repo.lastLoginUpdatedID != 99 {
			t.Fatalf("expected last login updated for id=99, got %d", repo.lastLoginUpdatedID)
		}
		if repo.lastLoginUpdatedTS < before || repo.lastLoginUpdatedTS > after {
			t.Fatalf("expected last login ts within [%d, %d], got %d", before, after, repo.lastLoginUpdatedTS)
		}
		if res.User.LastLoginAt != repo.lastLoginUpdatedTS {
			t.Fatalf("expected user.LastLoginAt updated to %d, got %d", repo.lastLoginUpdatedTS, res.User.LastLoginAt)
		}
		if res.User.PasswordHash != "" {
			t.Fatalf("expected returned user password hash to be scrubbed")
		}
		if res.Token == "" {
			t.Fatalf("expected token")
		}

		claims := parseHS256Claims(t, res.Token, secret)
		if claims["email"] != "a@b.com" {
			t.Fatalf("expected email claim, got %#v", claims["email"])
		}
		if claims["user_id"] == nil {
			t.Fatalf("expected user_id claim")
		}
		if claims["exp"] == nil {
			t.Fatalf("expected exp claim")
		}
	})

	t.Run("wrong password is rejected", func(t *testing.T) {
		t.Parallel()

		user := &model.User{
			BaseModel:    model.BaseModel{ID: 99},
			Email:        "a@b.com",
			PasswordHash: mustHash(t, "password123"),
			Status:       "active",
		}

		repo := &stubUserRepo{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return user, nil
			},
		}
		svc := NewAuthService(repo, []byte(secret), tokenTTL)

		_, err := svc.Login(context.Background(), LoginInput{
			Email:    "a@b.com",
			Password: "wrong",
		})
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Fatalf("expected ErrInvalidCredentials, got %v", err)
		}
		if repo.updateLastLogin != 0 {
			t.Fatalf("expected UpdateLastLogin not called, got %d", repo.updateLastLogin)
		}
	})

	t.Run("invalid token config is rejected", func(t *testing.T) {
		t.Parallel()

		repo := &stubUserRepo{}
		svc := NewAuthService(repo, nil, 0)

		_, err := svc.Login(context.Background(), LoginInput{
			Email:    "a@b.com",
			Password: "password123",
		})
		if !errors.Is(err, ErrInvalidTokenConfig) {
			t.Fatalf("expected ErrInvalidTokenConfig, got %v", err)
		}
	})
}
