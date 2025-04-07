package errs

import (
	"errors"
	"fmt"
)

var (
	ErrReadRequestData    = errors.New("failed to read request body")
	ErrParseRequestData   = errors.New("failed to parse request body")
	ErrNotFound           = errors.New("not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidID          = errors.New("invalid id format")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrBusinessLogic      = errors.New("business logic error")
	ErrProductNotFound         = errors.New("product not found")
	ErrProductNotApproved      = errors.New("product not approved")
	ErrNotEnoughStock          = errors.New("not enough stock")
	ErrProductDiscountNotFound = errors.New("product discount not found")
)

func NewBusinessLogicError(msg string) error {
	return fmt.Errorf("%w: %s", ErrBusinessLogic, msg)
}

func NewNotFoundError(msg string) error {
	return fmt.Errorf("%w: %s", ErrNotFound, msg)
}

func NewAlreadyExistsError(msg string) error {
	return fmt.Errorf("%w: %s", ErrAlreadyExists, msg)
}
