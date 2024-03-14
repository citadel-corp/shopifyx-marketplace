package bankaccount

import (
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
	"github.com/google/uuid"
)

type BankAccount struct {
	ID                uint64
	UUID              uuid.UUID
	BankName          string
	BankAccountName   string
	BankAccountNumber string
	User              user.User
}
