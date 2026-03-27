package service

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"time"

	"go-oj/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user disabled")
	ErrInvalidTokenConfig = errors.New("invalid token config")
)

const (
	UserStatusActive = "active"
)

type AuthServiceAPI interface {
	Register(ctx context.Context, input RegisterInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
}

type userRepo interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateLastLogin(ctx context.Context, id uint, ts int64) error
}

type AuthService struct {
	users     userRepo
	jwtSecret []byte
	tokenTTL  time.Duration
	initErr   error
}

func NewAuthService(users userRepo, jwtSecret []byte, tokenTTL time.Duration) *AuthService {
	service := &AuthService{
		users:     users,
		jwtSecret: bytes.TrimSpace(jwtSecret),
		tokenTTL:  tokenTTL,
	}
	if len(service.jwtSecret) == 0 || service.tokenTTL <= 0 {
		service.initErr = ErrInvalidTokenConfig
	}

	return service
}

type RegisterInput struct {
	Email    string
	Password string
	Nickname string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResult, error) {
	if s.initErr != nil {
		return nil, s.initErr
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	nickname := strings.TrimSpace(input.Nickname)
	password := input.Password

	if !isValidEmail(email) || !isValidPassword(password) || nickname == "" {
		return nil, ErrInvalidInput
	}

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrDuplicateEmail
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Nickname:     nickname,
		Status:       UserStatusActive,
		LastLoginAt:  0,
	}
	if err := s.users.Create(ctx, u); err != nil {
		if isDuplicateEmailError(err) {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}

	token, err := s.signToken(u)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: sanitizeUser(u)}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	if s.initErr != nil {
		return nil, s.initErr
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := input.Password

	if !isValidEmail(email) || password == "" {
		return nil, ErrInvalidInput
	}

	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !isUserActive(u.Status) {
		return nil, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	ts := time.Now().Unix()
	if err := s.users.UpdateLastLogin(ctx, u.ID, ts); err != nil {
		return nil, err
	}
	u.LastLoginAt = ts

	token, err := s.signToken(u)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: sanitizeUser(u)}, nil
}

func (s *AuthService) signToken(u *model.User) (string, error) {
	if s.initErr != nil {
		return "", s.initErr
	}

	exp := time.Now().Add(s.tokenTTL).Unix()
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"exp":     exp,
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tk.SignedString(s.jwtSecret)
}

func isValidEmail(email string) bool {
	// MVP validation: enough to reject obvious garbage.
	if email == "" {
		return false
	}
	if strings.Contains(email, " ") {
		return false
	}
	at := strings.IndexByte(email, '@')
	if at <= 0 || at == len(email)-1 {
		return false
	}
	if strings.Count(email, "@") != 1 {
		return false
	}
	// Require at least one dot after '@'.
	dot := strings.LastIndexByte(email, '.')
	return dot > at+1 && dot < len(email)-1
}

func isValidPassword(pw string) bool {
	// MVP: length check only.
	return len(pw) >= 6
}

func isDuplicateEmailError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key") || strings.Contains(message, "unique constraint")
}

func isUserActive(status string) bool {
	return status == "" || status == UserStatusActive
}

func sanitizeUser(user *model.User) *model.User {
	if user == nil {
		return nil
	}

	safeUser := *user
	safeUser.PasswordHash = ""
	return &safeUser
}
