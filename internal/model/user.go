package model

type User struct {
	BaseModel
	Email        string `gorm:"size:128;not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Nickname     string `gorm:"size:64;not null" json:"nickname"`
	Status       string `gorm:"size:32;not null;default:active;index" json:"status"`
	LastLoginAt  int64  `gorm:"not null;default:0" json:"last_login_at"`
}
