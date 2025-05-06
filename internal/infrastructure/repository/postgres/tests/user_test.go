package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	user "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
)

func TestUserRepository_GetUserByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		email := "test@example.com"
		expectedUser := &models.UserDB{
			ID:           userID,
			Email:        email,
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hash"),
			ImageURL:     null.StringFrom("image.jpg"),
			Role:         models.RoleBuyer,
		}

		row := sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "password_hash", "image_url", "role",
		}).AddRow(
			expectedUser.ID,
			expectedUser.Email,
			expectedUser.Name,
			expectedUser.Surname,
			expectedUser.PasswordHash,
			expectedUser.ImageURL,
			expectedUser.Role.String(),
		)

		mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
			WithArgs(email).
			WillReturnRows(row)

		user, err := repo.GetUserByEmail(context.Background(), email)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		email := "notfound@example.com"

		mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(context.Background(), email)
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrInvalidCredentials))
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		email := "error@example.com"

		mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
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

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := &models.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hash"),
			ImageURL:     null.StringFrom("image.jpg"),
			PhoneNumber:  null.StringFrom("+1234567890"),
			Role:         models.RoleBuyer,
		}

		row := sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "password_hash", "image_url", "phone_number", "role",
		}).AddRow(
			expectedUser.ID,
			expectedUser.Email,
			expectedUser.Name,
			expectedUser.Surname,
			expectedUser.PasswordHash,
			expectedUser.ImageURL,
			expectedUser.PhoneNumber,
			expectedUser.Role.String(),
		)

		mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
			WithArgs(userID).
			WillReturnRows(row)

		user, err := repo.GetUserByID(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("not found", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByID(context.Background(), userID)
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrInvalidCredentials))
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
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

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "new_image.jpg"

		mock.ExpectExec(`UPDATE bazaar.user SET image_url = \$1 WHERE id = \$2`).
			WithArgs(imageURL, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserImageURL(context.Background(), userID, imageURL)
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "new_image.jpg"

		mock.ExpectExec(`UPDATE bazaar.user SET image_url = \$1 WHERE id = \$2`).
			WithArgs(imageURL, userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserImageURL(context.Background(), userID, imageURL)
		require.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrNotFound))
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "new_image.jpg"

		mock.ExpectExec(`UPDATE bazaar.user SET image_url = \$1 WHERE id = \$2`).
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

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		updateData := models.UpdateUserDB{
			Name:        "NewName",
			Surname:     null.StringFrom("NewSurname"),
			PhoneNumber: null.StringFrom("+1234567890"),
		}

		mock.ExpectExec(`UPDATE bazaar.user SET name = \$1, surname = \$2, phone_number = \$3 WHERE id = \$4`).
			WithArgs(updateData.Name, updateData.Surname, updateData.PhoneNumber, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserProfile(context.Background(), userID, updateData)
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		updateData := models.UpdateUserDB{
			Name:        "NewName",
			Surname:     null.StringFrom("NewSurname"),
			PhoneNumber: null.StringFrom("+1234567890"),
		}

		mock.ExpectExec(`UPDATE bazaar.user SET name = \$1, surname = \$2, phone_number = \$3 WHERE id = \$4`).
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

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		newEmail := "new@example.com"

		mock.ExpectExec(`UPDATE bazaar.user SET email = \$1 WHERE id = \$2`).
			WithArgs(newEmail, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserEmail(context.Background(), userID, newEmail)
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		newEmail := "new@example.com"

		mock.ExpectExec(`UPDATE bazaar.user SET email = \$1 WHERE id = \$2`).
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

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		passwordHash := []byte("new_hash")

		mock.ExpectExec(`UPDATE bazaar.user SET password_hash = \$1 WHERE id = \$2`).
			WithArgs(passwordHash, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUserPassword(context.Background(), userID, passwordHash)
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		passwordHash := []byte("new_hash")

		mock.ExpectExec(`UPDATE bazaar.user SET password_hash = \$1 WHERE id = \$2`).
			WithArgs(passwordHash, userID).
			WillReturnError(errors.New("database error"))

		err := repo.UpdateUserPassword(context.Background(), userID, passwordHash)
		require.Error(t, err)
	})
}

func TestUserRepository_CreateSellerAndUpdateRole(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := user.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		title := "Test Seller"
		description := "Test Description"

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO bazaar.seller`).
			WithArgs(sqlmock.AnyArg(), title, description, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`UPDATE bazaar.user SET role = \$1 WHERE id = \$2`).
			WithArgs(models.RolePending.String(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.CreateSellerAndUpdateRole(context.Background(), userID, title, description)
		require.NoError(t, err)
	})

	t.Run("transaction begin error", func(t *testing.T) {
		userID := uuid.New()
		title := "Test Seller"
		description := "Test Description"

		mock.ExpectBegin().WillReturnError(errors.New("begin error"))

		err := repo.CreateSellerAndUpdateRole(context.Background(), userID, title, description)
		require.Error(t, err)
	})

	t.Run("create seller error", func(t *testing.T) {
		userID := uuid.New()
		title := "Test Seller"
		description := "Test Description"

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO bazaar.seller`).
			WithArgs(sqlmock.AnyArg(), title, description, userID).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		err := repo.CreateSellerAndUpdateRole(context.Background(), userID, title, description)
		require.Error(t, err)
	})

	t.Run("update role error", func(t *testing.T) {
		userID := uuid.New()
		title := "Test Seller"
		description := "Test Description"

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO bazaar.seller`).
			WithArgs(sqlmock.AnyArg(), title, description, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`UPDATE bazaar.user SET role = \$1 WHERE id = \$2`).
			WithArgs(models.RolePending.String(), userID).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := repo.CreateSellerAndUpdateRole(context.Background(), userID, title, description)
		require.Error(t, err)
	})

	t.Run("commit error", func(t *testing.T) {
		userID := uuid.New()
		title := "Test Seller"
		description := "Test Description"

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO bazaar.seller`).
			WithArgs(sqlmock.AnyArg(), title, description, userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`UPDATE bazaar.user SET role = \$1 WHERE id = \$2`).
			WithArgs(models.RolePending.String(), userID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit().WillReturnError(errors.New("commit error"))

		err := repo.CreateSellerAndUpdateRole(context.Background(), userID, title, description)
		require.Error(t, err)
	})
}