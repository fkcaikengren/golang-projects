package model

type ProblemTag struct {
	BaseModel
	ProblemID uint `gorm:"not null;uniqueIndex:uk_problem_tag;index" json:"problem_id"`
	TagID     uint `gorm:"not null;uniqueIndex:uk_problem_tag;index" json:"tag_id"`
}
