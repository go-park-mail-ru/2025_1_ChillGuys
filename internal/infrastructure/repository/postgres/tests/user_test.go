package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	userRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_GetUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()
		now := time.Now()

		expectedUser := &models.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hashed_password"),
			ImageURL:     null.StringFrom("image.jpg"),
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: now,
			},
		}

		row := sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "password_hash", "image_url",
			"user_version_id", "version", "updated_at",
		}).
			AddRow(
				expectedUser.ID,
				expectedUser.Email,
				expectedUser.Name,
				expectedUser.Surname,
				expectedUser.PasswordHash,
				expectedUser.ImageURL,
				expectedUser.UserVersion.ID,
				expectedUser.UserVersion.Version,
				expectedUser.UserVersion.UpdatedAt,
			)

		mock.ExpectQuery(`
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
      FROM bazaar.user u
        LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
      WHERE u.email = \$1;
    `).
			WithArgs(expectedUser.Email).
			WillReturnRows(row)

		user, err := repo.GetUserByEmail(context.Background(), expectedUser.Email)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		email := "notfound@example.com"

		mock.ExpectQuery(`
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
      FROM bazaar.user u
        LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
      WHERE u.email = \$1;
    `).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(context.Background(), email)
		require.Error(t, err)
		require.ErrorIs(t, err, errs.ErrNotFound)
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		email := "test@example.com"

		mock.ExpectQuery(`
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
      FROM bazaar.user u
        LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
      WHERE u.email = \$1;
    `).
			WithArgs(email).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByEmail(context.Background(), email)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()
		now := time.Now()

		expectedUser := &models.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hashed_password"),
			ImageURL:     null.StringFrom("image.jpg"),
			PhoneNumber:  null.StringFrom("+1234567890"),
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: now,
			},
		}

		row := sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "password_hash", "image_url", "phone_number",
			"user_version_id", "version", "updated_at",
		}).
			AddRow(
				expectedUser.ID,
				expectedUser.Email,
				expectedUser.Name,
				expectedUser.Surname,
				expectedUser.PasswordHash,
				expectedUser.ImageURL,
				expectedUser.PhoneNumber,
				expectedUser.UserVersion.ID,
				expectedUser.UserVersion.Version,
				expectedUser.UserVersion.UpdatedAt,
			)

		mock.ExpectQuery(`
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
      FROM bazaar.user u
      LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
      WHERE u.id = \$1;
    `).
			WithArgs(userID).
			WillReturnRows(row)

		user, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`
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
      FROM bazaar.user u
      LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
      WHERE u.id = \$1;
    `).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(context.Background(), userID)
		require.Error(t, err)
		require.ErrorIs(t, err, errs.ErrNotFound)
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`
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
      FROM bazaar.user u
      LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
      WHERE u.id = \$1;
    `).
			WithArgs(userID).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByID(context.Background(), userID)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_UpdateUserImageURL(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "new_image.jpg"

		mock.ExpectExec(`
      UPDATE bazaar.user SET image_url = \$1 WHERE id = \$2
    `).
			WithArgs(imageURL, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserImageURL(context.Background(), userID, imageURL)
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "new_image.jpg"

		mock.ExpectExec(`
      UPDATE bazaar.user SET image_url = \$1 WHERE id = \$2
    `).
			WithArgs(imageURL, userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserImageURL(context.Background(), userID, imageURL)
		require.Error(t, err)
		require.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "new_image.jpg"

		mock.ExpectExec(`
      UPDATE bazaar.user SET image_url = \$1 WHERE id = \$2
    `).
			WithArgs(imageURL, userID).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateUserImageURL(context.Background(), userID, imageURL)
		require.Error(t, err)
	})
}

func TestUserRepository_UpdateUserProfile(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		updateData := models.UpdateUserDB{
			Name:        "New Name",
			Surname:     null.StringFrom("New Surname"),
			PhoneNumber: null.StringFrom("+1234567890"),
		}

		mock.ExpectExec(`
      UPDATE bazaar.user SET name = \$1, surname = \$2, phone_number = \$3 WHERE id = \$4;
    `).
			WithArgs(updateData.Name, updateData.Surname, updateData.PhoneNumber, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserProfile(context.Background(), userID, updateData)
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		updateData := models.UpdateUserDB{
			Name:        "New Name",
			Surname:     null.StringFrom("New Surname"),
			PhoneNumber: null.StringFrom("+1234567890"),
		}

		mock.ExpectExec(`
      UPDATE bazaar.user SET name = \$1, surname = \$2, phone_number = \$3 WHERE id = \$4;
    `).
			WithArgs(updateData.Name, updateData.Surname, updateData.PhoneNumber, userID).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateUserProfile(context.Background(), userID, updateData)
		require.Error(t, err)
	})
}

func TestUserRepository_UpdateUserEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		newEmail := "new@example.com"

		mock.ExpectExec(`
      UPDATE bazaar.user SET email = \$1 WHERE id = \$2;
    `).
			WithArgs(newEmail, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserEmail(context.Background(), userID, newEmail)
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		newEmail := "new@example.com"

		mock.ExpectExec(`
      UPDATE bazaar.user SET email = \$1 WHERE id = \$2;
    `).
			WithArgs(newEmail, userID).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateUserEmail(context.Background(), userID, newEmail)
		require.Error(t, err)
	})
}

func TestUserRepository_UpdateUserPassword(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		newPasswordHash := []byte("new_hashed_password")

		mock.ExpectExec(`
      UPDATE bazaar.user SET password_hash = \$1 WHERE id = \$2;
    `).
			WithArgs(newPasswordHash, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserPassword(context.Background(), userID, newPasswordHash)
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		newPasswordHash := []byte("new_hashed_password")

		mock.ExpectExec(`
      UPDATE bazaar.user SET password_hash = \$1 WHERE id = \$2;
    `).
			WithArgs(newPasswordHash, userID).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateUserPassword(context.Background(), userID, newPasswordHash)
		require.Error(t, err)
	})
}
