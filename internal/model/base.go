package model

import "gorm.io/gorm"

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type SystemSetting struct {
	BaseModel
	Key   string `gorm:"size:100;uniqueIndex;not null" json:"key"`
	Value string `gorm:"size:255;not null" json:"value"`
}
