package bankaccount

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/db"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, bankAccount *BankAccount) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*BankAccount, error)
	List(ctx context.Context, userID uint64) ([]*BankAccount, error)
	Update(ctx context.Context, bankAccount *BankAccount) error
	Delete(ctx context.Context, uuid uuid.UUID) error
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
	if err := row.Err(); err != nil {
		return err
	}
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

// GetByUUID implements Repository.
func (d *dbRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*BankAccount, error) {
	getUserQuery := `
		SELECT b.uid, b.name, b.account_name, b.account_number, u.id
		FROM bank_accounts b
		INNER JOIN users u on b.user_id = u.id
		WHERE uid = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, uuid)
	i := &BankAccount{}
	err := row.Scan(&i.UUID, &i.User.Name, &i.BankAccountName, &i.BankAccountNumber, &i.User.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return i, nil
}

// Delete implements Repository.
func (d *dbRepository) Delete(ctx context.Context, uid uuid.UUID) error {
	deleteQuery := `
		DELETE FROM bank_accounts
		WHERE uid = $1;
	`
	row, err := d.db.DB().ExecContext(ctx, deleteQuery, uid)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// List implements Repository.
func (d *dbRepository) List(ctx context.Context, userID uint64) ([]*BankAccount, error) {
	listQuery := `
		SELECT uid, name, account_name, account_number
		FROM bank_accounts
		WHERE user_id = $1;
	`
	rows, err := d.db.DB().QueryContext(ctx, listQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var bankAccounts []*BankAccount
	for rows.Next() {
		i := &BankAccount{}
		if err := rows.Scan(&i.UUID, &i.BankName, &i.BankAccountName, &i.BankAccountNumber); err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bankAccounts, nil
}

// Update implements Repository.
func (d *dbRepository) Update(ctx context.Context, bankAccount *BankAccount) error {
	updateQuery := `
		UPDATE bank_accounts
		SET name = $1,
		account_name = $2,
		account_number = $3
		WHERE uid = $4;
	`
	_, err := d.db.DB().ExecContext(ctx, updateQuery, bankAccount.BankName, bankAccount.BankAccountName, bankAccount.BankAccountNumber, bankAccount.UUID)
	return err
}
