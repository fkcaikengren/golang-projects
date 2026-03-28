package service

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"time"

	"go-oj/internal/model"
	"go-oj/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 定义常量
const (
	AdminUserStatusActive   = "active"
	AdminUserStatusDisabled = "disabled"
	adminTokenScope         = "admin"
)

var ErrInvalidToken = errors.New("invalid token")

type AdminAuthServiceAPI interface {
	Login(ctx context.Context, input AdminLoginInput) (*AdminAuthResult, error)
	Authenticate(ctx context.Context, bearerToken string) (*model.AdminUser, error)
}

type adminUserRepo = repository.AdminUserRepository

type AdminAuthService struct {
	users     repository.AdminUserRepository
	jwtSecret []byte
	tokenTTL  time.Duration
	initErr   error
}

type AdminLoginInput struct {
	Email    string
	Password string
}

type AdminAuthResult struct {
	Token     string           `json:"token"`
	AdminUser *model.AdminUser `json:"admin_user"`
}

func NewAdminAuthService(users repository.AdminUserRepository, jwtSecret []byte, tokenTTL time.Duration) *AdminAuthService {
	service := &AdminAuthService{
		users:     users,
		jwtSecret: bytes.TrimSpace(jwtSecret),
		tokenTTL:  tokenTTL,
	}
	if len(service.jwtSecret) == 0 || service.tokenTTL <= 0 {
		service.initErr = ErrInvalidTokenConfig
	}

	return service
}

func (s *AdminAuthService) Login(ctx context.Context, input AdminLoginInput) (*AdminAuthResult, error) {
	if s.initErr != nil {
		return nil, s.initErr
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := input.Password
	if !isValidEmail(email) || password == "" {
		return nil, ErrInvalidInput
	}

	adminUser, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if !isAdminUserActive(adminUser.Status) {
		return nil, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	ts := time.Now().Unix()
	if err := s.users.UpdateLastLogin(ctx, adminUser.ID, ts); err != nil {
		return nil, err
	}
	adminUser.LastLoginAt = ts

	token, err := s.signToken(adminUser)
	if err != nil {
		return nil, err
	}

	return &AdminAuthResult{
		Token:     token,
		AdminUser: sanitizeAdminUser(adminUser),
	}, nil
}

func (s *AdminAuthService) Authenticate(ctx context.Context, bearerToken string) (*model.AdminUser, error) {
	if s.initErr != nil {
		return nil, s.initErr
	}

	token := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(bearerToken), "Bearer "))
	if token == "" || token == bearerToken {
		return nil, ErrInvalidToken
	}

	parsed, err := jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		if tk.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	scope, _ := claims["scope"].(string)
	if scope != adminTokenScope {
		return nil, ErrInvalidToken
	}

	adminUserID, err := uintFromClaim(claims["admin_user_id"])
	if err != nil {
		return nil, ErrInvalidToken
	}

	adminUser, err := s.users.GetByID(ctx, adminUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	if !isAdminUserActive(adminUser.Status) {
		return nil, ErrUserDisabled
	}

	return sanitizeAdminUser(adminUser), nil
}

func (s *AdminAuthService) signToken(user *model.AdminUser) (string, error) {
	if s.initErr != nil {
		return "", s.initErr
	}

	exp := time.Now().Add(s.tokenTTL).Unix()
	claims := jwt.MapClaims{
		"admin_user_id": user.ID,
		"email":         user.Email,
		"scope":         adminTokenScope,
		"exp":           exp,
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tk.SignedString(s.jwtSecret)
}

func isAdminUserActive(status string) bool {
	return status == "" || status == AdminUserStatusActive
}

func sanitizeAdminUser(user *model.AdminUser) *model.AdminUser {
	if user == nil {
		return nil
	}

	safeUser := *user
	safeUser.PasswordHash = ""
	return &safeUser
}

func uintFromClaim(value interface{}) (uint, error) {
	switch typed := value.(type) {
	case float64:
		return uint(typed), nil
	case int64:
		return uint(typed), nil
	case int:
		return uint(typed), nil
	default:
		return 0, ErrInvalidToken
	}
}
