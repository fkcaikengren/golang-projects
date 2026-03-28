package model

type ProblemSetProblem struct {
	BaseModel
	ProblemSetID uint `gorm:"not null;uniqueIndex:uk_set_problem;index;foreignKey:ProblemSetID;references:ID;constraint:OnDelete:CASCADE" json:"problem_set_id"`
	ProblemID    uint `gorm:"not null;uniqueIndex:uk_set_problem;index;foreignKey:ProblemID;references:ID;constraint:OnDelete:CASCADE" json:"problem_id"`
	SortOrder    int  `gorm:"not null;default:0" json:"sort_order"`
}
