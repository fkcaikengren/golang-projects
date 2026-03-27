package bootstrap

import (
	"strings"
	"testing"

	"go-oj/internal/config"
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
