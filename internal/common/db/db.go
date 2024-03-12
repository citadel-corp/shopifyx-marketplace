package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	sqlDB *sql.DB
}

func Connect(dbURL string) (*DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DB{sqlDB: db}, nil
}

func (db *DB) DB() *sql.DB {
	return db.sqlDB
}

func (db *DB) StartTx(ctx context.Context, f func(*sql.Tx) error) error {
	tx, err := db.sqlDB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) UpMigration() error {
	m, err := db.createMigrate()
	if err != nil {
		return err
	}

	limit := 0
	for {
		if err = m.Up(); err != nil {
			err = db.migrationErrorHandler(err, m)
		}
		if err == nil {
			slog.Info("Successfully running up migrations.")
			return nil
		}

		limit += 1
		if limit == 5 {
			slog.Error("Failed running up migrations.")
			return err
		}
	}
}

func (db *DB) DownMigration() error {
	m, err := db.createMigrate()
	if err != nil {
		return err
	}

	limit := 0
	for {
		if err = m.Down(); err != nil {
			err = db.migrationErrorHandler(err, m)
		}
		if err == nil {
			slog.Info("Successfully running down migrations.")
			return nil
		}

		limit += 1
		if limit == 5 {
			slog.Error("Failed running down migrations.")
			return err
		}
	}
}

func (db *DB) createMigrate() (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.sqlDB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		os.Getenv("MIGRATIONS_URI"),
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (db *DB) migrationErrorHandler(err error, m *migrate.Migrate) error {
	slog.Info(fmt.Sprintf("migration: handling error: %v", err))
	if strings.Contains(err.Error(), "Dirty database") {
		re := regexp.MustCompile("[0-9]+")
		s := re.FindAllString(err.Error(), -1)

		if len(s) > 0 {
			version, _ := strconv.Atoi(s[0])
			err := m.Force(version)
			if err != nil {
				return err
			}
		}
	} else if err == migrate.ErrNoChange {
		return nil
	}
	return err
}
