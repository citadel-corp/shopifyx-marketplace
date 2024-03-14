package bankaccount

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrNotFound         = errors.New("bank account not found")
	ErrUnauthorized     = errors.New("you are unauthorized to make changes to this bank account")
)
