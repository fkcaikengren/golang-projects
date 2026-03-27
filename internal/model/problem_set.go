package model

type ProblemSet struct {
	BaseModel
	Name        string `gorm:"size:128;not null" json:"name"`
	Slug        string `gorm:"size:128;not null;uniqueIndex" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"size:32;not null;default:active;index" json:"status"`
	SortOrder   int    `gorm:"not null;default:0" json:"sort_order"`
}
