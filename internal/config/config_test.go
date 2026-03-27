package config

import (
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	t.Setenv("APP_NAME", "go-oj")
	t.Setenv("APP_ENV", "test")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "oj")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "Postgres")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	t.Setenv("AUTH_TOKEN_TTL", "120")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppName != "go-oj" {
		t.Fatalf("unexpected AppName: %s", cfg.AppName)
	}

	if cfg.Database.Port != 5432 {
		t.Fatalf("unexpected DB_PORT: %d", cfg.Database.Port)
	}

	if cfg.Auth.JWTSecret != "test-secret" {
		t.Fatalf("unexpected AUTH_JWT_SECRET: %s", cfg.Auth.JWTSecret)
	}

	if cfg.Auth.TokenTTL != 120 {
		t.Fatalf("unexpected AUTH_TOKEN_TTL: %d", cfg.Auth.TokenTTL)
	}
}

func TestValidateMissingRequiredField(t *testing.T) {
	t.Setenv("APP_NAME", "go-oj")
	t.Setenv("APP_ENV", "test")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "oj")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("AUTH_JWT_SECRET", "test-secret")
	t.Setenv("AUTH_TOKEN_TTL", "120")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DB_PASSWORD is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateMissingAuthSecret(t *testing.T) {
	t.Setenv("APP_NAME", "go-oj")
	t.Setenv("APP_ENV", "test")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "oj")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "Postgres")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("AUTH_JWT_SECRET", "")
	t.Setenv("AUTH_TOKEN_TTL", "120")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}

	if !strings.Contains(err.Error(), "AUTH_JWT_SECRET is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}
