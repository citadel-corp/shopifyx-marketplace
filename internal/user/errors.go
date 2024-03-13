package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrWrongPassword         = errors.New("wrong password")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrValidationFailed      = errors.New("validation failed")
)
