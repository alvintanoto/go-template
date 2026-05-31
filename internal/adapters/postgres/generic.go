package postgres

import (
	"context"

	"gorm.io/gorm"
)

type GormRepository[DomainEntity any, GormModel any, ID any] struct {
	DB           *gorm.DB
	ToGorm       func(*DomainEntity) *GormModel
	ToDomain     func(*GormModel) *DomainEntity
	IDColumnName string
}

func (r *GormRepository[DomainEntity, GormModel, ID]) Create(ctx context.Context, entity *DomainEntity) error {
	dbModel := r.ToGorm(entity)
	return r.DB.WithContext(ctx).Create(dbModel).Error
}

func (r *GormRepository[DomainEntity, GormModel, ID]) GetByID(ctx context.Context, id ID) (*DomainEntity, error) {
	var dbModel GormModel
	err := r.DB.WithContext(ctx).First(&dbModel, r.IDColumnName+" = ?", id).Error
	if err != nil {
		return nil, err
	}
	return r.ToDomain(&dbModel), nil
}

func (r *GormRepository[DomainEntity, GormModel, ID]) Update(ctx context.Context, entity *DomainEntity) error {
	dbModel := r.ToGorm(entity)
	return r.DB.WithContext(ctx).Save(dbModel).Error
}

func (r *GormRepository[DomainEntity, GormModel, ID]) Delete(ctx context.Context, id ID) error {
	var dbModel GormModel
	return r.DB.WithContext(ctx).Where(r.IDColumnName+" = ?", id).Delete(&dbModel).Error
}
