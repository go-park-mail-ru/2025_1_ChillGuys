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
	ErrProductNotApproved = errors.New("product not approved")
	ErrNotEnoughStock     = errors.New("not enough stock")

	ErrMissingToken      = errors.New("missing jwt token")
	ErrNoMetadata        = errors.New("metadata is not provided")
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")
	ErrInternal          = errors.New("internal server error")
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
