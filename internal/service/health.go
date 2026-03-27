package service

import (
	"context"
	"time"

	"go-oj/internal/repository"
)

type HealthService interface {
	Health(ctx context.Context) map[string]interface{}
	Ready(ctx context.Context) error
}

type healthService struct {
	appName string
	repo    repository.HealthRepository
}

func NewHealthService(appName string, repo repository.HealthRepository) HealthService {
	return &healthService{
		appName: appName,
		repo:    repo,
	}
}

func (s *healthService) Health(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"app":       s.appName,
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *healthService) Ready(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
