package errs

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidUserID      = errors.New("invalid user id format")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
