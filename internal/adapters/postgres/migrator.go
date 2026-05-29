package postgres

import (
	"database/sql"
	"errors"

	"alvintanoto.id/go-template/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

func RunAutoMigrations(sqlDB *sql.DB, log *zap.Logger) error {
	sourceDriver, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("Database schema is up to date. No migrations applied.")
			return nil
		}
		return err
	}

	log.Info("DB migrations success")
	return nil
}
