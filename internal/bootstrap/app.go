package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-oj/internal/config"
	"go-oj/internal/handler"
	"go-oj/internal/model"
	"go-oj/internal/repository"
	"go-oj/internal/router"
	"go-oj/internal/service"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db, err := InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	healthRepo := repository.NewHealthRepository(db)
	healthService := service.NewHealthService(cfg.AppName, healthRepo)
	healthHandler := handler.NewHealthHandler(healthService)

	userRepo := repository.NewUserRepository(db)
	adminUserRepo := repository.NewAdminUserRepository(db)
	problemSetRepo := repository.NewProblemSetRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)

	tokenTTL := time.Duration(cfg.Auth.TokenTTL) * time.Minute
	authService := service.NewAuthService(userRepo, []byte(cfg.Auth.JWTSecret), tokenTTL)
	adminAuthService := service.NewAdminAuthService(adminUserRepo, []byte(cfg.Auth.JWTSecret), tokenTTL)
	adminAuthorizer, err := service.NewAdminAuthorizer(db)
	if err != nil {
		return nil, fmt.Errorf("new admin authorizer: %w", err)
	}
	if err := adminAuthorizer.SeedPolicies(context.Background()); err != nil {
		return nil, fmt.Errorf("seed admin policies: %w", err)
	}
	bootstrapAdmin, err := ensureBootstrapAdmin(context.Background(), adminUserRepo, cfg)
	if err != nil {
		return nil, fmt.Errorf("ensure bootstrap admin: %w", err)
	}
	if err := adminAuthorizer.AssignRole(strconv.FormatUint(uint64(bootstrapAdmin.ID), 10), service.AdminRoleAdmin); err != nil {
		return nil, fmt.Errorf("assign bootstrap admin role: %w", err)
	}
	problemSetService := service.NewProblemSetService(problemSetRepo)
	problemService := service.NewProblemService(problemRepo)
	submissionService := service.NewSubmissionService(submissionRepo, problemRepo)

	authHandler := handler.NewAuthHandler(authService)
	adminAuthHandler := handler.NewAdminAuthHandler(adminAuthService)
	adminHandler := handler.NewAdminHandler()
	problemSetHandler := handler.NewProblemSetHandler(problemSetService)
	problemHandler := handler.NewProblemHandler(problemService)
	submissionHandler := handler.NewSubmissionHandler(submissionService)

	return &App{
		Config: cfg,
		DB:     db,
		Router: router.New(
			healthHandler,
			authHandler,
			problemSetHandler,
			problemHandler,
			submissionHandler,
			adminAuthHandler,
			adminAuthService,
			adminHandler,
			adminAuthorizer,
		),
	}, nil
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.SystemSetting{},
		&model.User{},
		&model.AdminUser{},
		&model.ProblemSet{},
		&model.Problem{},
		&model.Tag{},
		&model.ProblemTag{},
		&model.ProblemSetProblem{},
		&model.Submission{},
		&model.UserProblemStat{},
		&model.OperationLog{},
	); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	log.Printf("database connected to %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return db, nil
}

type bootstrapAdminUserRepo interface {
	Create(ctx context.Context, user *model.AdminUser) error
	GetByEmail(ctx context.Context, email string) (*model.AdminUser, error)
}

func ensureBootstrapAdmin(ctx context.Context, repo bootstrapAdminUserRepo, cfg *config.Config) (*model.AdminUser, error) {
	email := strings.TrimSpace(strings.ToLower(cfg.AdminBootstrap.Email))
	if email == "" {
		email = "admin@go-oj.dev"
	}

	existing, err := repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	password := strings.TrimSpace(cfg.AdminBootstrap.Password)
	if password == "" {
		password = "Admin@123456"
	}

	displayName := strings.TrimSpace(cfg.AdminBootstrap.DisplayName)
	if displayName == "" {
		displayName = "Super Admin"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	adminUser := &model.AdminUser{
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Status:       service.AdminUserStatusActive,
		LastLoginAt:  0,
	}
	if err := repo.Create(ctx, adminUser); err != nil {
		return nil, err
	}

	return adminUser, nil
}
