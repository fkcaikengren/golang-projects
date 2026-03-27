package repository

import (
	"context"
	"strings"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

type ProblemFilter struct {
	Difficulty string
	TagSlug    string
	Keyword    string
}

type ProblemRepository interface {
	List(ctx context.Context, filter ProblemFilter) ([]model.Problem, error)
	GetBySlug(ctx context.Context, slug string) (*model.Problem, error)
	GetByID(ctx context.Context, id uint) (*model.Problem, error)
	ListTags(ctx context.Context, problemID uint) ([]model.Tag, error)
}

type problemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) ProblemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) List(ctx context.Context, filter ProblemFilter) ([]model.Problem, error) {
	var problems []model.Problem

	tx := r.db.WithContext(ctx).Model(&model.Problem{})
	if difficulty := strings.TrimSpace(filter.Difficulty); difficulty != "" {
		tx = tx.Where("problems.difficulty = ?", difficulty)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		tx = tx.Where("problems.title ILIKE ? OR problems.description ILIKE ?", like, like)
	}
	if tagSlug := strings.TrimSpace(filter.TagSlug); tagSlug != "" {
		tx = tx.
			Joins("JOIN problem_tags pt ON pt.problem_id = problems.id").
			Joins("JOIN tags ON tags.id = pt.tag_id").
			Where("tags.slug = ?", tagSlug)
	}

	err := tx.
		Distinct("problems.*").
		Order("problems.id ASC").
		Find(&problems).Error
	return problems, err
}

func (r *problemRepository) GetBySlug(ctx context.Context, slug string) (*model.Problem, error) {
	var problem model.Problem
	if err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&problem).Error; err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *problemRepository) GetByID(ctx context.Context, id uint) (*model.Problem, error) {
	var problem model.Problem
	if err := r.db.WithContext(ctx).First(&problem, id).Error; err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *problemRepository) ListTags(ctx context.Context, problemID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.WithContext(ctx).
		Model(&model.Tag{}).
		Joins("JOIN problem_tags pt ON pt.tag_id = tags.id").
		Where("pt.problem_id = ?", problemID).
		Order("tags.id ASC").
		Find(&tags).Error
	return tags, err
}
