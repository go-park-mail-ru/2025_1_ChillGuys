package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	userRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()
		now := time.Now()

		user := dto.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hashedpassword"),
			ImageURL:     null.String{},
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: now,
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO bazaar.\"user\"").
			WithArgs(user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.ImageURL).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.\"user_version\"").
			WithArgs(user.UserVersion.ID, user.UserVersion.UserID, user.UserVersion.Version, user.UserVersion.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		// For the basket creation, we need to use sqlmock.AnyArg since we can't predict the UUID
		mock.ExpectExec("INSERT INTO bazaar.basket").
			WithArgs(sqlmock.AnyArg(), user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), user)
		require.NoError(t, err)
	})

	t.Run("error in user insert", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()

		user := dto.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			PasswordHash: []byte("hashedpassword"),
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: time.Now(),
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO bazaar.\"user\"").
			WithArgs(user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.ImageURL).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), user)
		require.Error(t, err)
	})

	t.Run("error in version insert", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()

		user := dto.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			PasswordHash: []byte("hashedpassword"),
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: time.Now(),
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO bazaar.\"user\"").
			WithArgs(user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.ImageURL).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.\"user_version\"").
			WithArgs(user.UserVersion.ID, user.UserVersion.UserID, user.UserVersion.Version, user.UserVersion.UpdatedAt).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), user)
		require.Error(t, err)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()
		email := "test@example.com"
		now := time.Now()

		expectedUser := &dto.UserDB{
			ID:           userID,
			Email:        email,
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hashedpassword"),
			ImageURL:     null.StringFrom("http://example.com/image.jpg"),
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: now,
			},
		}

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
			FROM bazaar."user" u
					 LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
			WHERE u.email = \$1;
		`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{
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
			))

		user, err := repo.GetByEmail(context.Background(), email)
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
			FROM bazaar."user" u
					 LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
			WHERE u.email = \$1;
		`).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByEmail(context.Background(), email)
		require.ErrorIs(t, err, errs.ErrNotFound)
		assert.Nil(t, user)
	})

	t.Run("database error", func(t *testing.T) {
		email := "error@example.com"

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
			FROM bazaar."user" u
					 LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
			WHERE u.email = \$1;
		`).
			WithArgs(email).
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByEmail(context.Background(), email)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		versionID := uuid.New()
		now := time.Now()

		expectedUser := &dto.UserDB{
			ID:           userID,
			Email:        "test@example.com",
			Name:         "Test",
			Surname:      null.StringFrom("User"),
			PasswordHash: []byte("hashedpassword"),
			ImageURL:     null.StringFrom("http://example.com/image.jpg"),
			UserVersion: models.UserVersionDB{
				ID:        versionID,
				UserID:    userID,
				Version:   1,
				UpdatedAt: now,
			},
		}

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
			FROM bazaar."user" u
			LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
			WHERE u.id = \$1;
		`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{
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
			))

		user, err := repo.GetByID(context.Background(), userID)
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
				uv.id AS user_version_id, 
				uv.version, 
				uv.updated_at
			FROM bazaar."user" u
			LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
			WHERE u.id = \$1;
		`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByID(context.Background(), userID)
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
				uv.id AS user_version_id, 
				uv.version, 
				uv.updated_at
			FROM bazaar."user" u
			LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id
			WHERE u.id = \$1;
		`).
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		user, err := repo.GetByID(context.Background(), userID)
		require.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_CheckExistence(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("exists", func(t *testing.T) {
		email := "exists@example.com"

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM bazaar."user" WHERE email = \$1\)`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.CheckExistence(context.Background(), email)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("not exists", func(t *testing.T) {
		email := "notexists@example.com"

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM bazaar."user" WHERE email = \$1\)`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.CheckExistence(context.Background(), email)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("database error", func(t *testing.T) {
		email := "error@example.com"

		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM bazaar."user" WHERE email = \$1\)`).
			WithArgs(email).
			WillReturnError(errors.New("db error"))

		exists, err := repo.CheckExistence(context.Background(), email)
		require.Error(t, err)
		assert.False(t, exists)
	})
}

func TestUserRepository_IncrementVersion(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := "123"

		mock.ExpectExec(`UPDATE bazaar."user_version" SET version = version \+ 1 WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.IncrementVersion(context.Background(), userID)
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		userID := "123"

		mock.ExpectExec(`UPDATE bazaar."user_version" SET version = version \+ 1 WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.IncrementVersion(context.Background(), userID)
		require.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		userID := "123"

		mock.ExpectExec(`UPDATE bazaar."user_version" SET version = version \+ 1 WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		err := repo.IncrementVersion(context.Background(), userID)
		require.Error(t, err)
	})
}

func TestUserRepository_CheckVersion(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("version matches", func(t *testing.T) {
		userID := "123"
		version := 5

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(version))

		result := repo.CheckVersion(context.Background(), userID, version)
		assert.True(t, result)
	})

	t.Run("version doesn't match", func(t *testing.T) {
		userID := "123"
		version := 5

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(version + 1))

		result := repo.CheckVersion(context.Background(), userID, version)
		assert.False(t, result)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := "123"
		version := 5

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		result := repo.CheckVersion(context.Background(), userID, version)
		assert.False(t, result)
	})

	t.Run("database error", func(t *testing.T) {
		userID := "123"
		version := 5

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		result := repo.CheckVersion(context.Background(), userID, version)
		assert.False(t, result)
	})
}

func TestUserRepository_UpdateImageURL(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "http://example.com/new-image.jpg"

		mock.ExpectExec(`UPDATE bazaar."user" SET image_url = \$1 WHERE id = \$2`).
			WithArgs(imageURL, userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateImageURL(context.Background(), userID, imageURL)
		require.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "http://example.com/new-image.jpg"

		mock.ExpectExec(`UPDATE bazaar."user" SET image_url = \$1 WHERE id = \$2`).
			WithArgs(imageURL, userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateImageURL(context.Background(), userID, imageURL)
		require.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()
		imageURL := "http://example.com/new-image.jpg"

		mock.ExpectExec(`UPDATE bazaar."user" SET image_url = \$1 WHERE id = \$2`).
			WithArgs(imageURL, userID).
			WillReturnError(errors.New("db error"))

		err := repo.UpdateImageURL(context.Background(), userID, imageURL)
		require.Error(t, err)
	})
}

func TestUserRepository_GetCurrentVersion(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := "123"
		version := 5

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(version))

		result, err := repo.GetCurrentVersion(context.Background(), userID)
		require.NoError(t, err)
		assert.Equal(t, version, result)
	})

	t.Run("not found", func(t *testing.T) {
		userID := "123"

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetCurrentVersion(context.Background(), userID)
		require.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		userID := "123"

		mock.ExpectQuery(`SELECT version FROM bazaar."user_version" WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		_, err := repo.GetCurrentVersion(context.Background(), userID)
		require.Error(t, err)
	})
}