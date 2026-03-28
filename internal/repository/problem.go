package repository

import (
	"context"
	"strings"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

type ProblemFilter struct {
	Page     int
	PageSize int
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

	query := r.db.WithContext(ctx).Model(&model.Problem{})

	// 预加载 Tags，避免 N+1 问题
	query = query.Preload("Tags")

	if difficulty := strings.TrimSpace(filter.Difficulty); difficulty != "" {
		query = query.Where("difficulty = ?", difficulty)
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
	if tagSlug := strings.TrimSpace(filter.TagSlug); tagSlug != "" {
		query = query.
			Joins("JOIN problem_tags pt ON pt.problem_id = problems.id").
			Joins("JOIN tags ON tags.id = pt.tag_id").
			Where("tags.slug = ?", tagSlug)
	}

	// 分页
	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.
		Distinct("problems.id").
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
