package postgres

import (
	"context"

	"alvintanoto.id/go-template/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{}
}

// Create implements user.Repository.
func (u *UserRepository) Create(ctx context.Context, entity *user.User) error {
	return u.db.WithContext(ctx).Create(entity).Error
}

// Delete implements user.Repository.
func (u *UserRepository) Delete(ctx context.Context, id int64) error {
	return u.db.WithContext(ctx).Delete(&user.User{}, id).Error
}

// GetByID implements user.Repository.
func (u *UserRepository) GetByID(ctx context.Context, id int64) (*user.User, error) {
	var dest user.User
	if err := u.db.WithContext(ctx).First(&dest, id).Error; err != nil {
		return nil, err
	}
	return &dest, nil
}

// Update implements user.Repository.
func (u *UserRepository) Update(ctx context.Context, entity *user.User) error {
	return u.db.WithContext(ctx).Save(entity).Error
}
