package tests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

const (
	queryGetUserByEmail = `
	SELECT
		u.id,
		u.email,
		u.name,
		u.surname,
		u.password_hash,
		u.image_url,
		u.role,
		uv.id AS user_version_id,
		uv.version,
		uv.updated_at
	FROM bazaar.user u
	LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
	WHERE u.email = $1;
	`
	queryGetUserByID = `
	SELECT 
		u.id, 
		u.email, 
		u.name, 
		u.surname, 
		u.password_hash, 
		u.image_url, 
		u.phone_number,
		u.role,
		uv.id AS user_version_id,
		uv.version,
		uv.updated_at
	FROM bazaar.user u
	LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
	WHERE u.id = $1;
	`
	queryUpdateUserImageURL = `UPDATE bazaar.user SET image_url = $1 WHERE id = $2`
	queryUpdateUser         = `UPDATE bazaar.user SET name = $1, surname = $2, phone_number = $3 WHERE id = $4;`
	queryUpdateUserPassword = `UPDATE bazaar.user SET password_hash = $1 WHERE id = $2;`
	queryUpdateUserEmail    = `UPDATE bazaar.user SET email = $1 WHERE id = $2;`

	queryCreateSeller = `
        INSERT INTO bazaar.seller (id, title, description, user_id)
        VALUES ($1, $2, $3, $4)`

    queryUpdateUserRole = `
        UPDATE bazaar.user
        SET role = $1
        WHERE id = $2`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) UpdateUserImageURL(ctx context.Context, userID uuid.UUID, imageURL string) error {
	res, err := r.db.ExecContext(ctx, queryUpdateUserImageURL, imageURL, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID uuid.UUID, in models.UpdateUserDB) error {
	res, err := r.db.ExecContext(ctx, queryUpdateUser,
		in.Name,
		in.Surname,
		in.PhoneNumber,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error {
	res, err := r.db.ExecContext(ctx, queryUpdateUserEmail,
		email,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	res, err := r.db.ExecContext(ctx, queryUpdateUserPassword,
		passwordHash,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserDB, error) {
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
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserDB, error) {
	var user models.UserDB
	var phoneNumber null.String

	err := r.db.QueryRowContext(ctx, queryGetUserByID, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.ImageURL,
		&phoneNumber,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	user.PhoneNumber = phoneNumber

	return &user, nil
}

func (r *UserRepository) CreateSellerAndUpdateRole(ctx context.Context, userID uuid.UUID, title, description string) error {
    const op = "UserRepository.CreateSellerAndUpdateRole"
    
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
    }
    
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    _, err = tx.ExecContext(ctx, queryCreateSeller, 
        uuid.New(), 
        title, 
        description, 
        userID,
    )
    if err != nil {
        return fmt.Errorf("%s: failed to create seller: %w", op, err)
    }
    
    _, err = tx.ExecContext(ctx, queryUpdateUserRole, models.RolePending.String(), userID)
    if err != nil {
        return fmt.Errorf("%s: failed to update user role: %w", op, err)
    }
    
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
    }
    
    return nil
}