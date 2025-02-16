package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Blxssy/AvitoTest/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	driverName        = "pgx"
	applicationSchema = "public"
)

func New(cfg config.PostgresConfig) (*sqlx.DB, error) {
	sqlxDB, err := sqlx.Open(driverName, cfg.DataSource)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	if err := sqlxDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return sqlxDB, nil
}

func RunMigrations(instance *sql.DB, cfg config.PostgresConfig) (uint, error) {
	if _, err := instance.Exec("create schema if not exists " + applicationSchema); err != nil {
		return 0, fmt.Errorf("create schema: %w", err)
	}

	driver, err := postgres.WithInstance(instance, &postgres.Config{
		SchemaName: applicationSchema,
	})
	if err != nil {
		return 0, fmt.Errorf("create driver with instance: %w", err)
	}

	migrateInst, err := migrate.NewWithDatabaseInstance("file://"+cfg.PathToMigrations, driverName, driver)
	if err != nil {
		return 0, fmt.Errorf("create migrate instance: %w", err)
	}

	if err = migrateInst.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return 0, fmt.Errorf("up migrations: %w", err)
	}

	version, _, err := migrateInst.Version()
	if err != nil {
		return 0, fmt.Errorf("migrate: %w", err)
	}
	return version, nil
}
