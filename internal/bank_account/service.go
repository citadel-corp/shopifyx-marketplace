package bankaccount

import (
	"context"
	"fmt"

	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req CreateUpdateBankAccountPayload, userID uint64) (*BankAccountResponse, error)
	List(ctx context.Context, userID uint64) ([]*BankAccountResponse, error)
	PartialUpdate(ctx context.Context, req CreateUpdateBankAccountPayload, uuid uuid.UUID, userID uint64) (*BankAccountResponse, error)
	Delete(ctx context.Context, uuid uuid.UUID, userID uint64) error
}

type bankAccountService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &bankAccountService{repository: repository}
}

func (s *bankAccountService) Create(ctx context.Context, req CreateUpdateBankAccountPayload, userID uint64) (*BankAccountResponse, error) {
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

// Delete implements Service.
func (s *bankAccountService) Delete(ctx context.Context, uuid uuid.UUID, userID uint64) error {
	bankAccount, err := s.repository.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if bankAccount.User.ID != userID {
		return ErrForbidden
	}
	return s.repository.Delete(ctx, uuid)
}

// List implements Service.
func (s *bankAccountService) List(ctx context.Context, userID uint64) ([]*BankAccountResponse, error) {
	bankAccounts, err := s.repository.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	resp := make([]*BankAccountResponse, len(bankAccounts))
	for i, bankAccount := range bankAccounts {
		resp[i] = &BankAccountResponse{
			BankAccountID:     bankAccount.UUID.String(),
			BankName:          bankAccount.BankName,
			BankAccountName:   bankAccount.BankAccountName,
			BankAccountNumber: bankAccount.BankAccountNumber,
		}
	}
	return resp, nil
}

// PartialUpdate implements Service.
func (s *bankAccountService) PartialUpdate(ctx context.Context, req CreateUpdateBankAccountPayload, uuid uuid.UUID, userID uint64) (*BankAccountResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	bankAccount, err := s.repository.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if bankAccount.User.ID != userID {
		return nil, ErrForbidden
	}
	s.applyPartialUpdate(req, bankAccount)
	err = s.repository.Update(ctx, bankAccount)
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

func (s *bankAccountService) applyPartialUpdate(req CreateUpdateBankAccountPayload, bankAccount *BankAccount) {
	if req.BankName != "" {
		bankAccount.BankName = req.BankName
	}
	if req.BankAccountName != "" {
		bankAccount.BankAccountName = req.BankAccountName
	}
	if req.BankAccountNumber != "" {
		bankAccount.BankAccountNumber = req.BankAccountNumber
	}
}
