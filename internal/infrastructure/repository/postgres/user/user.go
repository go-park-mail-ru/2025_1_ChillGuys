package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
		uv.id AS user_version_id,
		uv.version,
		uv.updated_at
	FROM bazaar."user" u
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
		uv.id AS user_version_id, 
		uv.version, 
		uv.updated_at
	FROM bazaar."user" u
	LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
	WHERE u.id = $1;
	`
	queryUpdateUserImageURL = `UPDATE bazaar."user" SET image_url = $1 WHERE id = $2`
	queryUpdateUser         = `UPDATE bazaar."user" SET name = $1, surname = $2, phone_number = $3 WHERE id = $4;`
	queryUpdateUserPassword = `UPDATE bazaar."user" SET password_hash = $1 WHERE id = $2;`
	queryUpdateUserEmail    = `UPDATE bazaar."user" SET email = $1 WHERE id = $2;`
)

type UserRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewUserRepository(db *sql.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
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

func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID uuid.UUID, in dto.UpdateUserDB) error {
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*dto.UserDB, error) {
	var user dto.UserDB

	if err := r.db.QueryRowContext(ctx, queryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.ImageURL,
		&user.UserVersion.ID,
		&user.UserVersion.Version,
		&user.UserVersion.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	user.UserVersion.UserID = user.ID

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*dto.UserDB, error) {
	var user dto.UserDB

	err := r.db.QueryRowContext(ctx, queryGetUserByID, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.ImageURL,
		&user.PhoneNumber,
		&user.UserVersion.ID,
		&user.UserVersion.Version,
		&user.UserVersion.UpdatedAt,
	)
	user.UserVersion.UserID = user.ID

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}
