package repository

import (
	"context"

	"go-oj/internal/model"

	"gorm.io/gorm"
)

type AdminUserRepository interface {
	Create(ctx context.Context, user *model.AdminUser) error
	GetByEmail(ctx context.Context, email string) (*model.AdminUser, error)
	GetByID(ctx context.Context, id uint) (*model.AdminUser, error)
	UpdateLastLogin(ctx context.Context, id uint, ts int64) error
}

type adminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) AdminUserRepository {
	return &adminUserRepository{db: db}
}

func (r *adminUserRepository) Create(ctx context.Context, user *model.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *adminUserRepository) GetByEmail(ctx context.Context, email string) (*model.AdminUser, error) {
	var user model.AdminUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *adminUserRepository) GetByID(ctx context.Context, id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *adminUserRepository) UpdateLastLogin(ctx context.Context, id uint, ts int64) error {
	result := r.db.WithContext(ctx).
		Model(&model.AdminUser{}).
		Where("id = ?", id).
		Update("last_login_at", ts)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
