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
	queryCreateUser           = `INSERT INTO "user" (user_id, email, name, surname, password_hash, version) VALUES($1, $2, $3, $4, $5, $6);`
	queryGetUserVersion       = `SELECT version FROM "user" WHERE user_id = $1`
	queryGetUserByEmail       = `SELECT user_id, email, name, surname, password_hash, version FROM "user" WHERE email = $1`
	queryGetUserByID          = `SELECT user_id, email, name, surname, password_hash, version FROM "user" WHERE user_id = $1`
	queryIncrementUserVersion = `UPDATE "user" SET version = version + 1 WHERE user_id = $1`
	queryCheckUserExists      = `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`
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
	_, err := r.db.ExecContext(ctx, queryCreateUser,
		user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.Version,
	)

	return err
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

	err := r.db.QueryRowContext(ctx, queryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Surname,
		&user.PasswordHash,
		&user.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
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
		&user.Version,
	)

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

	err := r.db.QueryRowContext(ctx, queryGetUserVersion, userID).Scan(&currentVersion)

	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	if err != nil {
		return false
	}

	return currentVersion == version
}

func (r *UserRepository) CheckUserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, queryCheckUserExists, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
