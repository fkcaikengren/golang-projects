package model

type UserProblemStat struct {
	BaseModel
	UserID           uint   `gorm:"not null;uniqueIndex:uk_user_problem;index;foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user_id"`
	ProblemID        uint   `gorm:"not null;uniqueIndex:uk_user_problem;index;foreignKey:ProblemID;references:ID;constraint:OnDelete:CASCADE" json:"problem_id"`
	Status           string `gorm:"size:32;not null;default:attempted;index" json:"status"`
	SubmitCount      int    `gorm:"not null;default:0" json:"submit_count"`
	AcceptedCount    int    `gorm:"not null;default:0" json:"accepted_count"`
	FirstAcceptedAt  int64  `gorm:"not null;default:0" json:"first_accepted_at"`
	LastSubmitAt     int64  `gorm:"not null;default:0" json:"last_submit_at"`
	LastSubmissionID uint   `gorm:"not null;default:0" json:"last_submission_id"`
}
