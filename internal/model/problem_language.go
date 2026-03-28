package model

type ProblemLanguage struct {
	BaseModel
	ProblemID      uint   `gorm:"not null;uniqueIndex:uk_problem_language;index;foreignKey:ProblemID;references:ID;constraint:OnDelete:CASCADE" json:"problem_id"`
	Language       string `gorm:"size:32;not null;uniqueIndex:uk_problem_language;index" json:"language"`
	TimeLimitMS    int    `gorm:"not null;default:1000" json:"time_limit_ms"`
	MemoryLimitKB  int    `gorm:"not null;default:65536" json:"memory_limit_kb"`
	SourceTemplate string `gorm:"type:text" json:"source_template"`
	Status         string `gorm:"size:32;not null;default:enabled;index" json:"status"`
}
