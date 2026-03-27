package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	AppEnv   string
	HTTPPort string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !strings.Contains(err.Error(), "no such file or directory") && err != configFileNotFound {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	cfg := &Config{
		AppName:  v.GetString("APP_NAME"),
		AppEnv:   v.GetString("APP_ENV"),
		HTTPPort: v.GetString("HTTP_PORT"),
		Database: DatabaseConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			Name:     v.GetString("DB_NAME"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			SSLMode:  v.GetString("DB_SSLMODE"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "go-oj")
	v.SetDefault("APP_ENV", "dev")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("DB_HOST", "127.0.0.1")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_SSLMODE", "disable")
}

func (c *Config) Validate() error {
	if c.HTTPPort == "" {
		return fmt.Errorf("HTTP_PORT is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if c.Database.Port == 0 {
		return fmt.Errorf("DB_PORT is required")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if c.Database.SSLMode == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}

	return nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host,
		d.Port,
		d.User,
		d.Password,
		d.Name,
		d.SSLMode,
	)
}
