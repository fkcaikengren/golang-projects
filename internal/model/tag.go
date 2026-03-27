package model

type Tag struct {
	BaseModel
	Name  string `gorm:"size:64;not null" json:"name"`
	Slug  string `gorm:"size:64;not null;uniqueIndex" json:"slug"`
	Type  string `gorm:"size:32;not null;index" json:"type"`
	Color string `gorm:"size:32" json:"color"`
}
