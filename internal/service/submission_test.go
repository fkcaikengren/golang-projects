package service

import (
	"context"
	"errors"
	"testing"

	"go-oj/internal/model"
)

type stubSubmissionRepo struct {
	lastSaved              *model.Submission
	lastListUserID         uint
	lastListUserProblemID  uint
	saveWithStatsFn        func(ctx context.Context, submission *model.Submission) error
	listByUserFn           func(ctx context.Context, userID uint) ([]model.Submission, error)
	listByUserAndProblemFn func(ctx context.Context, userID uint, problemID uint) ([]model.Submission, error)
}

func (s *stubSubmissionRepo) Create(ctx context.Context, submission *model.Submission) error {
	return nil
}

func (s *stubSubmissionRepo) ListByUser(ctx context.Context, userID uint) ([]model.Submission, error) {
	s.lastListUserID = userID
	if s.listByUserFn != nil {
		return s.listByUserFn(ctx, userID)
	}
	return nil, nil
}

func (s *stubSubmissionRepo) ListByUserAndProblem(ctx context.Context, userID uint, problemID uint) ([]model.Submission, error) {
	s.lastListUserID = userID
	s.lastListUserProblemID = problemID
	if s.listByUserAndProblemFn != nil {
		return s.listByUserAndProblemFn(ctx, userID, problemID)
	}
	return nil, nil
}

func (s *stubSubmissionRepo) SaveWithStats(ctx context.Context, submission *model.Submission) error {
	s.lastSaved = submission
	if s.saveWithStatsFn != nil {
		return s.saveWithStatsFn(ctx, submission)
	}
	submission.ID = 1
	return nil
}

type stubProblemLookupRepo struct {
	getByIDFn   func(ctx context.Context, id uint) (*model.Problem, error)
	getBySlugFn func(ctx context.Context, slug string) (*model.Problem, error)
}

func (s stubProblemLookupRepo) GetByID(ctx context.Context, id uint) (*model.Problem, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (s stubProblemLookupRepo) GetBySlug(ctx context.Context, slug string) (*model.Problem, error) {
	if s.getBySlugFn != nil {
		return s.getBySlugFn(ctx, slug)
	}
	return nil, errors.New("not found")
}

func TestSubmissionServiceSubmitAccepted(t *testing.T) {
	t.Parallel()

	repo := &stubSubmissionRepo{}
	svc := NewSubmissionService(
		repo,
		stubProblemLookupRepo{
			getByIDFn: func(ctx context.Context, id uint) (*model.Problem, error) {
				return &model.Problem{BaseModel: model.BaseModel{ID: id}, Slug: "two-sum"}, nil
			},
		},
	)

	res, err := svc.Submit(context.Background(), SubmitCodeInput{
		UserID:    1,
		ProblemID: 2,
		Language:  "go",
		Code:      "func solve() int { return 1 }",
	})
	if err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if res == nil {
		t.Fatalf("expected submission result")
	}
	if repo.lastSaved == nil {
		t.Fatalf("expected submission persisted")
	}
	if repo.lastSaved.Status != SubmissionStatusAccepted {
		t.Fatalf("saved status = %q, want %q", repo.lastSaved.Status, SubmissionStatusAccepted)
	}
	if repo.lastSaved.Score != 100 {
		t.Fatalf("saved score = %d, want 100", repo.lastSaved.Score)
	}
	if repo.lastSaved.JudgedAt == 0 || repo.lastSaved.SubmitAt == 0 {
		t.Fatalf("expected submit and judged timestamps to be set")
	}
}

func TestSubmissionServiceSubmitWrongAnswer(t *testing.T) {
	t.Parallel()

	repo := &stubSubmissionRepo{}
	svc := NewSubmissionService(
		repo,
		stubProblemLookupRepo{
			getByIDFn: func(ctx context.Context, id uint) (*model.Problem, error) {
				return &model.Problem{BaseModel: model.BaseModel{ID: id}, Slug: "two-sum"}, nil
			},
		},
	)

	res, err := svc.Submit(context.Background(), SubmitCodeInput{
		UserID:    1,
		ProblemID: 2,
		Language:  "go",
		Code:      "fmt.Println(1)",
	})
	if err != nil {
		t.Fatalf("Submit error: %v", err)
	}
	if res.Status != SubmissionStatusWrongAnswer {
		t.Fatalf("status = %q, want %q", res.Status, SubmissionStatusWrongAnswer)
	}
	if repo.lastSaved.Score != 0 {
		t.Fatalf("saved score = %d, want 0", repo.lastSaved.Score)
	}
}

func TestSubmissionServiceListUserSubmissions(t *testing.T) {
	t.Parallel()

	want := []model.Submission{{BaseModel: model.BaseModel{ID: 10}}}
	repo := &stubSubmissionRepo{
		listByUserFn: func(ctx context.Context, userID uint) ([]model.Submission, error) {
			return want, nil
		},
	}
	svc := NewSubmissionService(repo, stubProblemLookupRepo{})

	got, err := svc.ListUserSubmissions(context.Background(), 8)
	if err != nil {
		t.Fatalf("ListUserSubmissions error: %v", err)
	}
	if repo.lastListUserID != 8 {
		t.Fatalf("list user id = %d, want 8", repo.lastListUserID)
	}
	if len(got) != 1 || got[0].ID != 10 {
		t.Fatalf("unexpected submissions: %#v", got)
	}
}
