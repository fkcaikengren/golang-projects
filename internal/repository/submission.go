package repository

import (
	"context"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	Create(ctx context.Context, submission *model.Submission) error
	ListByUser(ctx context.Context, userID uint) ([]model.Submission, error)
	ListByUserAndProblem(ctx context.Context, userID uint, problemID uint) ([]model.Submission, error)
	SaveWithStats(ctx context.Context, submission *model.Submission) error
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(ctx context.Context, submission *model.Submission) error {
	return r.db.WithContext(ctx).Create(submission).Error
}

func (r *submissionRepository) ListByUser(ctx context.Context, userID uint) ([]model.Submission, error) {
	var submissions []model.Submission
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("submit_at DESC, id DESC").
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) ListByUserAndProblem(ctx context.Context, userID uint, problemID uint) ([]model.Submission, error) {
	var submissions []model.Submission
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND problem_id = ?", userID, problemID).
		Order("submit_at DESC, id DESC").
		Find(&submissions).Error
	return submissions, err
}

func (r *submissionRepository) SaveWithStats(ctx context.Context, submission *model.Submission) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(submission).Error; err != nil {
			return err
		}

		var stat model.UserProblemStat
		err := tx.Where("user_id = ? AND problem_id = ?", submission.UserID, submission.ProblemID).
			First(&stat).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
			stat = model.UserProblemStat{
				UserID:    submission.UserID,
				ProblemID: submission.ProblemID,
			}
		}

		stat.SubmitCount++
		stat.LastSubmitAt = submission.SubmitAt
		stat.LastSubmissionID = submission.ID

		if submission.Status == "accepted" {
			stat.AcceptedCount++
			if stat.FirstAcceptedAt == 0 {
				stat.FirstAcceptedAt = submission.JudgedAt
			}
			stat.Status = "solved"
		} else if stat.AcceptedCount > 0 {
			stat.Status = "solved"
		} else {
			stat.Status = "attempted"
		}

		if stat.ID == 0 {
			return tx.Create(&stat).Error
		}
		return tx.Save(&stat).Error
	})
}
