package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByID(ctx context.Context, id uint64) (*User, error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, user *User) error {
	createUserQuery := `
		INSERT INTO users (
			username, name, hashed_password
		) VALUES (
			$1, $2, $3
		)
		RETURNING id;
	`
	row := d.db.DB().QueryRowContext(ctx, createUserQuery, user.Username, user.Name, user.HashedPassword)
	var id uint64
	err := row.Scan(&id)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUsernameAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	user.ID = id
	return nil
}

// GetByUsernameAndHashedPassword implements Repository.
func (d *dbRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	getUserQuery := `
		SELECT id, username, name, hashed_password FROM users
		WHERE username = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, username)
	u := &User{}
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (d *dbRepository) GetByID(ctx context.Context, id uint64) (*User, error) {
	getUserQuery := `
		SELECT id, username, name, product_sold_total, hashed_password FROM users
		WHERE id = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.ProductSoldTotal, &u.HashedPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
