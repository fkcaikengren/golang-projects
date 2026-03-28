package service

import (
	"context"
	"strings"
	"time"

	"go-oj/internal/model"
	"go-oj/internal/repository"
)

const (
	SubmissionStatusAccepted    = "accepted"
	SubmissionStatusWrongAnswer = "wrong_answer"
)

type submissionRepo = repository.SubmissionRepository

type problemLookupRepo = repository.ProblemRepository

type SubmissionServiceAPI interface {
	Submit(ctx context.Context, input SubmitCodeInput) (*model.Submission, error)
	ListUserSubmissions(ctx context.Context, userID uint) ([]model.Submission, error)
	ListProblemSubmissions(ctx context.Context, userID uint, problemSlug string) ([]model.Submission, error)
}

type SubmissionService struct {
	submissions repository.SubmissionRepository
	problems    repository.ProblemRepository
}

type SubmitCodeInput struct {
	UserID    uint
	ProblemID uint
	Language  string
	Code      string
}

func NewSubmissionService(submissions repository.SubmissionRepository, problems repository.ProblemRepository) *SubmissionService {
	return &SubmissionService{
		submissions: submissions,
		problems:    problems,
	}
}

func (s *SubmissionService) Submit(ctx context.Context, input SubmitCodeInput) (*model.Submission, error) {
	if input.UserID == 0 || input.ProblemID == 0 || strings.TrimSpace(input.Language) == "" || strings.TrimSpace(input.Code) == "" {
		return nil, ErrInvalidInput
	}

	if _, err := s.problems.GetByID(ctx, input.ProblemID); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	submission := &model.Submission{
		UserID:    input.UserID,
		ProblemID: input.ProblemID,
		Language:  strings.TrimSpace(input.Language),
		Code:      input.Code,
		SubmitAt:  now,
		JudgedAt:  now,
		RuntimeMS: 1,
		MemoryKB:  64,
	}

	if strings.Contains(input.Code, "return") {
		submission.Status = SubmissionStatusAccepted
		submission.Score = 100
	} else {
		submission.Status = SubmissionStatusWrongAnswer
		submission.Score = 0
	}

	if err := s.submissions.SaveWithStats(ctx, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

func (s *SubmissionService) ListUserSubmissions(ctx context.Context, userID uint) ([]model.Submission, error) {
	if userID == 0 {
		return nil, ErrInvalidInput
	}
	return s.submissions.ListByUser(ctx, userID)
}

func (s *SubmissionService) ListProblemSubmissions(ctx context.Context, userID uint, problemSlug string) ([]model.Submission, error) {
	if userID == 0 || strings.TrimSpace(problemSlug) == "" {
		return nil, ErrInvalidInput
	}

	problem, err := s.problems.GetBySlug(ctx, problemSlug)
	if err != nil {
		return nil, err
	}

	return s.submissions.ListByUserAndProblem(ctx, userID, problem.ID)
}
