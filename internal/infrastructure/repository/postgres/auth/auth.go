package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryCreateUser = `
		INSERT INTO bazaar.user (id, email, name, surname, password_hash, image_url, role) 
		VALUES($1, $2, $3, $4, $5, $6, $7);
	`

	queryGetUserByEmail = `
		SELECT id, email, name, surname, password_hash, image_url, role
		FROM bazaar.user 
		WHERE email = $1;
	`

	queryGetUserByID = `
		SELECT id, email, name, surname, password_hash, image_url, phone_number, role
		FROM bazaar.user 
		WHERE id = $1;
	`

	queryCheckUserExists = `
		SELECT EXISTS(SELECT 1 FROM bazaar.user WHERE email = $1);
	`

	queryCreateBasket = `
		INSERT INTO bazaar.basket (id, user_id, total_price, total_price_discount)
		SELECT $1, $2, 0, 0;
	`
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.UserDB) error {
	const op = "AuthRepository.CreateUser"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("begin transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	user.Role = models.RoleBuyer
	_, err = tx.ExecContext(ctx, queryCreateUser,
		user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.ImageURL, user.Role,
	)
	if err != nil {
		logger.WithError(err).Error("create user")
		tx.Rollback()
		return err
	}

	basketID := uuid.New()
	_, err = tx.ExecContext(ctx, queryCreateBasket, basketID, user.ID)
	if err != nil {
		logger.WithError(err).Error("create basket")
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("commit transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserDB, error) {
	const op = "AuthRepository.GetUserByEmail"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var user models.UserDB

	err := r.db.QueryRowContext(ctx, queryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.ImageURL,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("user not found by email")
			return nil, errs.ErrInvalidCredentials
		}
		logger.WithError(err).Error("get user by email")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserDB, error) {
	const op = "AuthRepository.GetUserByID"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var user models.UserDB

	err := r.db.QueryRowContext(ctx, queryGetUserByID, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.ImageURL,
		&user.PhoneNumber,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("user not found by ID")
			return nil, errs.ErrInvalidCredentials
		}
		logger.WithError(err).Error("get user by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *AuthRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	const op = "AuthRepository.CheckUserExists"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var exists bool

	if err := r.db.QueryRowContext(ctx, queryCheckUserExists, email).Scan(&exists); err != nil {
		logger.WithError(err).Error("failed to check user existence")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}
