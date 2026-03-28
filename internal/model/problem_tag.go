package model

type ProblemTag struct {
	BaseModel
	ProblemID uint `gorm:"not null;uniqueIndex:uk_problem_tag;index;foreignKey:ProblemID;references:ID;constraint:OnDelete:CASCADE" json:"problem_id"`
	TagID     uint `gorm:"not null;uniqueIndex:uk_problem_tag;index;foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE" json:"tag_id"`
}
