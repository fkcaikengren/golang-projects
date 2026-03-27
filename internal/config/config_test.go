package config

import (
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

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error, got nil")
	}
}
