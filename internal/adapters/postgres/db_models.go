package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"alvintanoto.id/go-template/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormUser struct {
	ID        string         `gorm:"primaryKey;column:id;type:varchar(40);"`
	Name      string         `gorm:"column:name;not null;type:varchar(100);"`
	Email     string         `gorm:"column:email;unique;not null;type:varchar(100);"`
	Password  string         `gorm:"column:password;not null; type:varchar(255);"`
	CreatedAt time.Time      `gorm:"column:created_at;not null;default:now()"`
	CreatedBy string         `gorm:"column:created_by;type:varchar(40)"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null;default:now()"`
	UpdatedBy string         `gorm:"column:updated_by;type:varchar(40)"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (GormUser) TableName() string {
	return "users"
}

func NewGormDatabase(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbPort)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Error("Failed to open raw database connection", zap.Error(err))
		return nil, err
	}

	// run db auto migration
	if err := RunAutoMigrations(sqlDB, log); err != nil {
		log.Fatal("Database migration failed, halting startup", zap.Error(err))
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		ConnPool: sqlDB,
	})
	if err != nil {
		log.Error("connecting to database: ", zap.Error(err))
	}

	return db, err
}
