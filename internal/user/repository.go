package user

import (
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
)

type Repository interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
}

type DBRepository struct {
	db *db.DB
}

// Create implements Repository.
func (d *DBRepository) Create(user *User) error {
	panic("unimplemented")
}

// GetByUsernameAndHashedPassword implements Repository.
func (d *DBRepository) GetByUsername(username string) (*User, error) {
	panic("unimplemented")
}
