package model

type Submission struct {
	BaseModel
	UserID      uint   `gorm:"not null;index;foreignKey:UserID;references:ID;constraint:OnDelete:SET NULL" json:"user_id"`
	ProblemID   uint   `gorm:"not null;index;foreignKey:ProblemID;references:ID;constraint:OnDelete:CASCADE" json:"problem_id"`
	Language    string `gorm:"size:32;not null;index" json:"language"`
	Code        string `gorm:"type:text;not null" json:"code"`
	Status      string `gorm:"size:32;not null;default:pending;index" json:"status"`
	Score       int    `gorm:"not null;default:0" json:"score"`
	RuntimeMS   int    `gorm:"not null;default:0" json:"runtime_ms"`
	MemoryKB    int    `gorm:"not null;default:0" json:"memory_kb"`
	SubmitAt    int64  `gorm:"not null;index" json:"submit_at"`
	JudgedAt    int64  `gorm:"not null;default:0" json:"judged_at"`
}
