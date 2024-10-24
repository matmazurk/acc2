package db

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var fs embed.FS

func migrateUp(db *sqlx.DB) error {
	driver, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)

	}

	dbDriver, err := sqlite.WithInstance(db.DB, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	migrations, err := migrate.NewWithInstance("iofs", driver, "sqlite", dbDriver)
	if err != nil {
		return fmt.Errorf("could not create migrations instance: %w", err)
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not perform migrations up: %w", err)
	}

	return nil
}
