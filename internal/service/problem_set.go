package service

import (
	"context"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

var ErrProblemSetNotFound = gorm.ErrRecordNotFound

type problemSetReader interface {
	ListActive(ctx context.Context) ([]model.ProblemSet, error)
	GetBySlug(ctx context.Context, slug string) (*model.ProblemSet, error)
	ListProblems(ctx context.Context, setID uint) ([]model.Problem, error)
}

type ProblemSetServiceAPI interface {
	List(ctx context.Context) ([]model.ProblemSet, error)
	Detail(ctx context.Context, slug string) (*ProblemSetDetail, error)
}

type ProblemSetDetail struct {
	ProblemSet *model.ProblemSet `json:"problem_set"`
	Problems   []model.Problem   `json:"problems"`
}

type ProblemSetService struct {
	repo problemSetReader
}

func NewProblemSetService(repo problemSetReader) *ProblemSetService {
	return &ProblemSetService{repo: repo}
}

func (s *ProblemSetService) List(ctx context.Context) ([]model.ProblemSet, error) {
	return s.repo.ListActive(ctx)
}

func (s *ProblemSetService) Detail(ctx context.Context, slug string) (*ProblemSetDetail, error) {
	set, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	problems, err := s.repo.ListProblems(ctx, set.ID)
	if err != nil {
		return nil, err
	}

	return &ProblemSetDetail{
		ProblemSet: set,
		Problems:   problems,
	}, nil
}
