package service

import (
	"context"

	"go-oj/internal/model"
	"go-oj/internal/repository"

	"gorm.io/gorm"
)

var ErrProblemNotFound = gorm.ErrRecordNotFound

type problemReader = repository.ProblemRepository

type ProblemServiceAPI interface {
	List(ctx context.Context, input ProblemListInput) ([]model.Problem, error)
	Detail(ctx context.Context, slug string) (*ProblemDetail, error)
}

type ProblemListInput struct {
	Difficulty string
	Tag        string
	Keyword    string
}

type ProblemDetail struct {
	Problem *model.Problem `json:"problem"`
	Tags    []model.Tag    `json:"tags"`
}

type ProblemService struct {
	repo repository.ProblemRepository
}

func NewProblemService(repo repository.ProblemRepository) *ProblemService {
	return &ProblemService{repo: repo}
}

func (s *ProblemService) List(ctx context.Context, input ProblemListInput) ([]model.Problem, error) {
	return s.repo.List(ctx, repository.ProblemFilter{
		Difficulty: input.Difficulty,
		TagSlug:    input.Tag,
		Keyword:    input.Keyword,
	})
}

func (s *ProblemService) Detail(ctx context.Context, slug string) (*ProblemDetail, error) {
	problem, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	tags, err := s.repo.ListTags(ctx, problem.ID)
	if err != nil {
		return nil, err
	}

	return &ProblemDetail{
		Problem: problem,
		Tags:    tags,
	}, nil
}
