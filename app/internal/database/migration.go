package database

import (
	"log/slog"

	"hotel.com/app/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(conn string, l *slog.Logger) error {
	d, err := iofs.New(sql.SqlFiles, "migrations")
	if err != nil {
		l.Error("failed to create migration source:", "err", err.Error())
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, conn)
	if err != nil {
		l.Error("failed to create migration instance", "err", err.Error())
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		l.Error("migration failed:", "err", err.Error())
		return err
	}
	if err == migrate.ErrNoChange {
		l.Info("no pending migrations")
	} else {
		l.Info("no pending migrations")
	}

	return nil
}
