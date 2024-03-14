package bankaccount

import (
	"context"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, bankAccount *BankAccount) error
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, bankAccount *BankAccount) error {
	createUserQuery := `
		INSERT INTO bank_accounts (
			name, account_name, account_number, user_id
		) VALUES (
			$1, $2, $3, $4
		)
		RETURNING id, uid;
	`
	row := d.db.DB().QueryRowContext(ctx, createUserQuery, bankAccount.BankName, bankAccount.BankAccountName, bankAccount.BankAccountNumber, bankAccount.User.ID)
	var id uint64
	var uuid uuid.UUID
	err := row.Scan(&id, &uuid)
	if err != nil {
		return err
	}
	bankAccount.ID = id
	bankAccount.UUID = uuid
	return nil
}
