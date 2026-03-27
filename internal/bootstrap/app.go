package bootstrap

import (
	"fmt"
	"log"
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

	return &App{
		Config: cfg,
		DB:     db,
		Router: router.New(healthHandler),
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
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.SystemSetting{},
		&model.User{},
		&model.ProblemSet{},
		&model.Problem{},
		&model.Tag{},
		&model.ProblemTag{},
		&model.ProblemSetProblem{},
		&model.Submission{},
		&model.UserProblemStat{},
	); err != nil {
		return nil, err
	}

	log.Printf("database connected to %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	return db, nil
}
