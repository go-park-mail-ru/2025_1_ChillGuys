package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
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
		u.role
	FROM bazaar.user u
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
		u.role
	FROM bazaar.user u
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
		return errs.NewNotFoundError("user not found")
	}

	return nil
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID uuid.UUID, in models.UpdateUserDB) error {
	_, err := r.db.ExecContext(ctx, queryUpdateUser,
		in.Name,
		in.Surname,
		in.PhoneNumber,
		userID,
	)
	return err
}

func (r *UserRepository) UpdateUserEmail(ctx context.Context, userID uuid.UUID, email string) error {
	_, err := r.db.ExecContext(ctx, queryUpdateUserEmail,
		email,
		userID,
	)
	return err
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash []byte) error {
	_, err := r.db.ExecContext(ctx, queryUpdateUserPassword,
		passwordHash,
		userID,
	)
	return err
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserDB, error) {
	var user models.UserDB

	if err := r.db.QueryRowContext(ctx, queryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.ImageURL,
		&user.Role,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserDB, error) {
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
			return nil, errs.ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateSellerAndUpdateRole(ctx context.Context, userID uuid.UUID, title, description string)  error {
    const op = "UserRepository.CreateSellerAndUpdateRole"
    
    // Начинаем транзакцию
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
    }
    
    // Отложенный rollback в случае ошибки
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
    
    // 2. Обновляем роль пользователя
    _, err = tx.ExecContext(ctx, queryUpdateUserRole, models.RolePending.String(), userID)
    if err != nil {
        return fmt.Errorf("%s: failed to update user role: %w", op, err)
    }
    
    // Фиксируем транзакцию
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
    }
    
    return nil
}