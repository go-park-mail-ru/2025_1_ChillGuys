package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryCreateUser        = `INSERT INTO bazaar."user" (id, email, name, surname, password_hash, image_url) VALUES($1, $2, $3, $4, $5, $6);`
	queryCreateUserVersion = `INSERT INTO bazaar."user_version" (id, user_id, version, updated_at) VALUES($1, $2, $3, $4);`
	queryGetUserVersion    = `SELECT version FROM bazaar."user_version" WHERE user_id = $1 LIMIT 1`
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
		uv.id AS user_version_id, 
		uv.version, 
		uv.updated_at
	FROM bazaar."user" u
	LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
	WHERE u.id = $1;
	`
	queryIncrementUserVersion = `UPDATE bazaar."user_version" SET version = version + 1 WHERE user_id = $1`
	queryCheckUserExists      = `SELECT EXISTS(SELECT 1 FROM bazaar."user" WHERE email = $1)`
	queryUpdateUserImageURL   = `UPDATE bazaar."user" SET image_url = $1 WHERE id = $2`
	queryCreateBasket 		  = `
			INSERT INTO bazaar.basket (id, user_id, total_price, total_price_discount)
			SELECT $1, $2, 0, 0
	`
)

type UserRepository struct {
	db  *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:  db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user dto.UserDB) error {
	const op = "UserRepository.Create"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("begin transaction")
		return err
	}

	_, err = tx.ExecContext(ctx, queryCreateUser,
		user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.ImageURL,
	)
	if err != nil {
		logger.WithError(err).Error("create user")
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryCreateUserVersion,
		user.UserVersion.ID, user.UserVersion.UserID, user.UserVersion.Version, user.UserVersion.UpdatedAt,
	)
	if err != nil {
		logger.WithError(err).Error("create user version")
		tx.Rollback()
		return err
	}

	basketID := uuid.New()
    _, err = tx.ExecContext(ctx, queryCreateBasket,
        basketID,
        user.ID, 
    )
    if err != nil {
        logger.WithError(err).Error("failed to create basket")
        return err
    }

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("commit transaction")
		return err
	}
	return nil
}

func (r *UserRepository) GetCurrentVersion(ctx context.Context, userID string) (int, error) {
	const op = "UserRepository.GetCurrentVersion"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var version int

	err := r.db.QueryRowContext(ctx, queryGetUserVersion, userID).Scan(&version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("user version not found")
			return 0, errs.NewNotFoundError(op)
		}
		logger.WithError(err).Error("get user version")
		return 0, err
	}

	return version, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*dto.UserDB, error) {
	const op = "UserRepository.GetByEmail"
	logger := logctx.GetLogger(ctx).WithField("op", op)

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
			logger.Warn("user not found by email")
			return nil, errs.NewNotFoundError(op)
		}
		logger.WithError(err).Error("get user by email")
		return nil, err
	}

	user.UserVersion.UserID = user.ID

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*dto.UserDB, error) {
	const op = "UserRepository.GetByID"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var user dto.UserDB

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
			logger.Warn("user not found by ID")
			return nil, errs.NewNotFoundError(op)
		}
		logger.WithError(err).Error("get user by ID")
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) IncrementVersion(ctx context.Context, userID string) error {
	const op = "UserRepository.IncrementVersion"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	res, err := r.db.ExecContext(ctx, queryIncrementUserVersion, userID)
	if err != nil {
		logger.WithError(err).Error("increment version")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("get rows affected")
		return err
	}
	
	if rowsAffected == 0 {
		logger.Warn("no rows affected when incrementing version")
		return errs.NewNotFoundError(op)
	}

	return nil
}

func (r *UserRepository) CheckVersion(ctx context.Context, userID string, version int) bool {
	const op = "UserRepository.CheckVersion"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var currentVersion int

	if err := r.db.QueryRowContext(ctx, queryGetUserVersion, userID).Scan(&currentVersion); err != nil {
		logger.WithError(err).Error("check version")
		return false
	}

	return currentVersion == version
}

func (r *UserRepository) CheckExistence(ctx context.Context, email string) (bool, error) {
	const op = "UserRepository.CheckExistence"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var exists bool

	if err := r.db.QueryRowContext(ctx, queryCheckUserExists, email).Scan(&exists); err != nil {
		logger.WithError(err).Error("failed to check user existence")
		return false, err
	}

	return exists, nil
}

func (r *UserRepository) UpdateImageURL(ctx context.Context, userID uuid.UUID, imageURL string) error {
	const op = "UserRepository.UpdateImageURL"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	res, err := r.db.ExecContext(ctx, queryUpdateUserImageURL, imageURL, userID)
	if err != nil {
		logger.WithError(err).Error("update image URL")
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("get rows affected")
		return err
	}

	if rowsAffected == 0 {
		logger.Warn("no rows affected when updating image URL")
		return errs.NewNotFoundError(op)
	}

	return nil
}
