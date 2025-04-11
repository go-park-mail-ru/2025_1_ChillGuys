package tests

import (
	"context"
	"database/sql"
	"errors"
	user2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := user2.NewAuthRepository(db, logrus.New())

	userID := uuid.New()
	userVersionID := uuid.New()
	user := dto.UserDB{
		ID:           userID,
		Email:        "test@example.com",
		Name:         "Test",
		Surname:      null.StringFrom("User"),
		PasswordHash: []byte("hashedpassword"),
		ImageURL:     null.String{},
		UserVersion: models.UserVersionDB{
			ID:        userVersionID,
			UserID:    userID,
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO bazaar.\"user\"").WithArgs(
		user.ID,
		user.Email,
		user.Name,
		user.Surname,
		user.PasswordHash,
		user.ImageURL,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO bazaar.\"user_version\"").WithArgs(
		user.UserVersion.ID,
		user.UserVersion.UserID,
		user.UserVersion.Version,
		user.UserVersion.UpdatedAt,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := user2.NewAuthRepository(db, logrus.New())

	email := "test@example.com"
	userID := uuid.New()
	userVersionID := uuid.New()

	mock.ExpectQuery(`SELECT u.id, u.email, u.name, u.surname, u.password_hash, u.image_url, uv.id AS user_version_id, uv.version, uv.updated_at FROM bazaar."user" u LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id WHERE u.email = \$1;`).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "password_hash", "image_url", "user_version_id", "version", "updated_at",
		}).
			AddRow(userID, email, "Test", "User", []byte("hashedpassword"), "http://example.com/image.jpg", userVersionID, 1, time.Now()))

	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "http://example.com/image.jpg", user.ImageURL.String)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := user.NewUserRepository(db, logrus.New())

	userID := uuid.New()
	userVersionID := uuid.New()
	updatedAt := time.Now()

	mock.ExpectQuery(`SELECT u.id, u.email, u.name, u.surname, u.password_hash, u.image_url, u.phone_number, uv.id AS user_version_id, uv.version, uv.updated_at FROM bazaar."user" u LEFT JOIN bazaar.user_version uv ON u.id = uv.user_id WHERE u.id = \$1;`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "password_hash", "image_url", "phone_number",
			"user_version_id", "version", "updated_at",
		}).
			AddRow(
				userID,
				"test@example.com",
				"Test",
				"User",
				[]byte("hashedpassword"),
				"http://example.com/image.jpg",
				"+79999999999", // <-- добавь фиктивный номер
				userVersionID,
				1,
				updatedAt,
			))

	user, err := repo.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "http://example.com/image.jpg", user.ImageURL.String)
	assert.Equal(t, userVersionID, user.UserVersion.ID)
	assert.Equal(t, 1, user.UserVersion.Version)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCheckUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := user2.NewAuthRepository(db, logrus.New())

	email := "test@example.com"

	mock.ExpectQuery("SELECT EXISTS").WithArgs(email).WillReturnRows(
		sqlmock.NewRows([]string{"exists"}).AddRow(true),
	)

	exists, err := repo.CheckUserExists(context.Background(), email)
	assert.NoError(t, err)
	assert.True(t, exists)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestIncrementUserVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := user2.NewAuthRepository(db, logrus.New())

	userID := "123"
	expectedQuery := `UPDATE bazaar."user_version" SET version = version \+ 1 WHERE user_id = \$1`

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.IncrementUserVersion(context.Background(), userID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.IncrementUserVersion(context.Background(), userID)
		assert.True(t, errors.Is(err, errs.ErrNotFound))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectExec(expectedQuery).
			WithArgs(userID).
			WillReturnError(sql.ErrConnDone)

		err := repo.IncrementUserVersion(context.Background(), userID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCheckUserVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := user2.NewAuthRepository(db, logrus.New())

	userID := "123"
	version := 5
	expectedQuery := `SELECT version FROM bazaar."user_version" WHERE user_id = \$1`

	t.Run("VersionMatches", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(version))

		result := repo.CheckUserVersion(context.Background(), userID, version)
		assert.True(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("VersionDoesNotMatch", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(version + 1))

		result := repo.CheckUserVersion(context.Background(), userID, version)
		assert.False(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		result := repo.CheckUserVersion(context.Background(), userID, version)
		assert.False(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnError(sql.ErrConnDone)

		result := repo.CheckUserVersion(context.Background(), userID, version)
		assert.False(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
