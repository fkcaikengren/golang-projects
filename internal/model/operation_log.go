package model

type OperationLog struct {
	BaseModel
	AdminUserID uint   `gorm:"not null;index" json:"admin_user_id"`
	Resource    string `gorm:"size:64;not null" json:"resource"`
	Action      string `gorm:"size:32;not null" json:"action"`
	TargetType  string `gorm:"size:64;not null" json:"target_type"`
	TargetID    string `gorm:"size:64;not null" json:"target_id"`
	RequestID   string `gorm:"size:64;not null" json:"request_id"`
	DetailJSON  string `gorm:"type:text;not null" json:"detail_json"`
	IP          string `gorm:"size:64;not null" json:"ip"`
	UserAgent   string `gorm:"size:255;not null" json:"user_agent"`
}
