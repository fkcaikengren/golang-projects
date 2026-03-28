package model

type JudgeConfig struct {
	BaseModel
	Name               string `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Language           string `gorm:"size:32;not null;index" json:"language"`
	CompileCommand     string `gorm:"type:text" json:"compile_command"`
	RunCommand         string `gorm:"type:text;not null" json:"run_command"`
	TimeLimitMS        int    `gorm:"not null;default:1000" json:"time_limit_ms"`
	MemoryLimitKB      int    `gorm:"not null;default:65536" json:"memory_limit_kb"`
	SpecialJudgeMode   string `gorm:"size:32;not null;default:none" json:"special_judge_mode"`
	SpecialJudgeSource string `gorm:"type:text" json:"special_judge_source"`
	SandboxConfigJSON  string `gorm:"type:text;not null;default:'{}'" json:"sandbox_config_json"`
	Status             string `gorm:"size:32;not null;default:active;index" json:"status"`
}
