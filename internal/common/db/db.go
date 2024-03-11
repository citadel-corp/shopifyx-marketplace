package db

import (
	"context"
	"database/sql"

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

func (db *DB) UpMigration() {
	// pake golang migrate disini
}

func (db *DB) DownMigration() {
	// pake golang migrate disini
}
