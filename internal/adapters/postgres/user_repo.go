package postgres

import (
	"context"

	"alvintanoto.id/go-template/internal/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	*GormRepository[user.User, GormUser, string]
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{
		db: db,
		GormRepository: &GormRepository[user.User, GormUser, string]{
			DB:           db,
			IDColumnName: "id",
			ToGorm: func(u *user.User) *GormUser {
				return &GormUser{ID: u.UserID, Email: u.Email, Password: u.Password}
			},
			ToDomain: func(m *GormUser) *user.User {
				return &user.User{UserID: m.ID, Email: m.Email, Password: m.Password}
			},
		},
	}
}

// FindByEmail implements user.Repository.
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var dbUser GormUser

	err := u.db.WithContext(ctx).
		Where("email = ?", email).
		First(&dbUser).
		Error

	if err != nil {
		return nil, err
	}

	return u.GormRepository.ToDomain(&dbUser), nil
}

// ExistsByEmail implements user.Repository.
func (u *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := u.db.WithContext(ctx).
		Table("users").
		Select("1").
		Where("email = ?", email).
		Limit(1).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}
