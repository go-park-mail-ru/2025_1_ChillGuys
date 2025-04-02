package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

const (
	queryCreateUser        = `INSERT INTO "user" (id, email, name, surname, password_hash, image_url) VALUES($1, $2, $3, $4, $5, $6);`
	queryCreateUserVersion = `INSERT INTO "user_version" (id, user_id, version, updated_at) VALUES($1, $2, $3, $4);`
	queryGetUserVersion    = `SELECT version FROM "user_version" WHERE user_id = $1`
	queryGetUserByEmail    = `
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
	FROM "user" u
			 LEFT JOIN user_version uv ON u.id = uv.user_id
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
		uv.id AS user_version_id, 
		uv.version, 
		uv.updated_at
	FROM "user" u
	LEFT JOIN user_version uv ON u.id = uv.user_id
	WHERE u.id = $1;
	`
	queryIncrementUserVersion = `UPDATE "user_version" SET version = version + 1 WHERE user_id = $1`
	queryCheckUserExists      = `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`
	queryUpdateUserImageURL   = `UPDATE "user" SET image_url = $1 WHERE id = $2`
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

func (r *UserRepository) CreateUser(ctx context.Context, user models.UserDB) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryCreateUser,
		user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.ImageURL,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryCreateUserVersion,
		user.UserVersion.ID, user.UserVersion.UserID, user.UserVersion.Version, user.UserVersion.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetUserCurrentVersion(ctx context.Context, userID string) (int, error) {
	var version int

	err := r.db.QueryRowContext(ctx, queryGetUserVersion, userID).Scan(&version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrUserNotFound
		}

		return 0, err
	}

	return version, nil
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
		&user.UserVersion.ID,
		&user.UserVersion.Version,
		&user.UserVersion.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	user.UserVersion.UserID = user.ID

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
		&user.UserVersion.ID,
		&user.UserVersion.Version,
		&user.UserVersion.UpdatedAt,
	)
	user.UserVersion.UserID = user.ID

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) IncrementUserVersion(ctx context.Context, userID string) error {
	res, err := r.db.ExecContext(ctx, queryIncrementUserVersion, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) CheckUserVersion(ctx context.Context, userID string, version int) bool {
	var currentVersion int

	if err := r.db.QueryRowContext(ctx, queryGetUserVersion, userID).Scan(&currentVersion); err != nil {
		return false
	}

	return currentVersion == version
}

func (r *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	if err := r.db.QueryRowContext(ctx, queryCheckUserExists, email).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
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
		return models.ErrUserNotFound
	}

	return nil
}
