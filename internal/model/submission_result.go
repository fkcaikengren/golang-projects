package model

type SubmissionResult struct {
	BaseModel
	SubmissionID  uint   `gorm:"not null;index;foreignKey:SubmissionID;references:ID;constraint:OnDelete:CASCADE" json:"submission_id"`
	TestCaseID    uint   `gorm:"not null;default:0;index;foreignKey:TestCaseID;references:ID;constraint:OnDelete:SET NULL" json:"test_case_id"`
	CaseOrder     int    `gorm:"not null;default:0" json:"case_order"`
	Status        string `gorm:"size:32;not null;index" json:"status"`
	RuntimeMS     int    `gorm:"not null;default:0" json:"runtime_ms"`
	MemoryKB      int    `gorm:"not null;default:0" json:"memory_kb"`
	ExitCode      int    `gorm:"not null;default:0" json:"exit_code"`
	CheckerLog    string `gorm:"type:text" json:"checker_log"`
	ErrorMessage  string `gorm:"type:text" json:"error_message"`
	InputSnapshot string `gorm:"type:text" json:"input_snapshot"`
	OutputActual  string `gorm:"type:text" json:"output_actual"`
	OutputExpect  string `gorm:"type:text" json:"output_expect"`
}
