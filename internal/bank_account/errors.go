package bankaccount

import "errors"

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrNotFound         = errors.New("bank account not found")
	ErrForbidden        = errors.New("you are forbidden to make changes to this bank account")
)
