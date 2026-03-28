package model

type TestCase struct {
	BaseModel
	ProblemID   uint   `gorm:"not null;index;foreignKey:ProblemID;references:ID;constraint:OnDelete:CASCADE" json:"problem_id"`
	Name        string `gorm:"size:128;not null" json:"name"`
	Input       string `gorm:"type:text;not null" json:"input"`
	Output      string `gorm:"type:text;not null" json:"output"`
	IsSample    bool   `gorm:"not null;default:false;index" json:"is_sample"`
	SortOrder   int    `gorm:"not null;default:0" json:"sort_order"`
	Status      string `gorm:"size:32;not null;default:active;index" json:"status"`
	InputMode   string `gorm:"size:32;not null;default:stdin" json:"input_mode"`
	OutputMode  string `gorm:"size:32;not null;default:stdout" json:"output_mode"`
	Explanation string `gorm:"type:text" json:"explanation"`
	Version     int    `gorm:"not null;default:1" json:"version"`
}
