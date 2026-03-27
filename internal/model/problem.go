package model

type Problem struct {
	BaseModel
	Title             string `gorm:"size:255;not null;index" json:"title"`
	Slug              string `gorm:"size:255;not null;uniqueIndex" json:"slug"`
	Source            string `gorm:"size:128" json:"source"`
	Difficulty        string `gorm:"size:32;not null;index" json:"difficulty"`
	Description       string `gorm:"type:text;not null" json:"description"`
	InputDescription  string `gorm:"type:text" json:"input_description"`
	OutputDescription string `gorm:"type:text" json:"output_description"`
	SampleInput       string `gorm:"type:text" json:"sample_input"`
	SampleOutput      string `gorm:"type:text" json:"sample_output"`
	Hint              string `gorm:"type:text" json:"hint"`
	Status            string `gorm:"size:32;not null;default:draft;index" json:"status"`
}
