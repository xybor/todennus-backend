package postgres

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./infras/database/postgres/migration",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
