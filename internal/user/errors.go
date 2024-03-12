package user

import "errors"

var (
	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
	ErrUsernameAlreadyExists   = errors.New("username already exists")
	ErrValidationFailed        = errors.New("validation failed")
)
