package user

import (
	"context"
	"errors"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByUsername(username string) (*User, error)
}

type DBRepository struct {
	db *db.DB
}

// Create implements Repository.
func (d *DBRepository) Create(ctx context.Context, user *User) error {
	createUserQuery := `
		INSERT INTO users (
			username, name, hashed_password
		) VALUES (
			$1, $2, $3
		)
		RETURNING id;
	`
	row := d.db.DB().QueryRowContext(ctx, createUserQuery, user.Username, user.Name, user.HashedPassword)
	var id int
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
func (d *DBRepository) GetByUsername(username string) (*User, error) {
	panic("unimplemented")
}
