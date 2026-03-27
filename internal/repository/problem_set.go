package repository

import (
	"context"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

type ProblemSetRepository interface {
	ListActive(ctx context.Context) ([]model.ProblemSet, error)
	GetBySlug(ctx context.Context, slug string) (*model.ProblemSet, error)
	ListProblems(ctx context.Context, setID uint) ([]model.Problem, error)
}

type problemSetRepository struct {
	db *gorm.DB
}

func NewProblemSetRepository(db *gorm.DB) ProblemSetRepository {
	return &problemSetRepository{db: db}
}

func (r *problemSetRepository) ListActive(ctx context.Context) ([]model.ProblemSet, error) {
	var sets []model.ProblemSet
	err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("sort_order ASC, id ASC").
		Find(&sets).Error
	return sets, err
}

func (r *problemSetRepository) GetBySlug(ctx context.Context, slug string) (*model.ProblemSet, error) {
	var set model.ProblemSet
	if err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&set).Error; err != nil {
		return nil, err
	}
	return &set, nil
}

func (r *problemSetRepository) ListProblems(ctx context.Context, setID uint) ([]model.Problem, error) {
	var problems []model.Problem
	err := r.db.WithContext(ctx).
		Model(&model.Problem{}).
		Joins("JOIN problem_set_problems psp ON psp.problem_id = problems.id").
		Where("psp.problem_set_id = ?", setID).
		Order("psp.sort_order ASC, problems.id ASC").
		Find(&problems).Error
	return problems, err
}
