package user

import (
	"context"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
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
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// GetByUsernameAndHashedPassword implements Repository.
func (d *DBRepository) GetByUsername(username string) (*User, error) {
	panic("unimplemented")
}
