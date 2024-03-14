package bankaccount

import (
	"context"
	"fmt"

	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type Service interface {
	Create(ctx context.Context, req CreateBankAccountPayload, userID uint64) (*BankAccountResponse, error)
}

type bankAccountService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &bankAccountService{repository: repository}
}

func (s *bankAccountService) Create(ctx context.Context, req CreateBankAccountPayload, userID uint64) (*BankAccountResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	bankAccount := &BankAccount{
		BankName:          req.BankName,
		BankAccountName:   req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber,
		User: user.User{
			ID: userID,
		},
	}

	err = s.repository.Create(ctx, bankAccount)
	if err != nil {
		return nil, err
	}
	return &BankAccountResponse{
		BankAccountID:     bankAccount.UUID.String(),
		BankName:          bankAccount.BankName,
		BankAccountName:   bankAccount.BankAccountName,
		BankAccountNumber: bankAccount.BankAccountNumber,
	}, nil
}
