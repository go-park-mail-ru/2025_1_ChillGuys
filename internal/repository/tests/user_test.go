package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/guregu/null"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db, logrus.New())

	user := models.UserDB{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Name:         "Test",
		Surname:      null.StringFrom("User"),
		PasswordHash: []byte("hashedpassword"),
		Version:      1,
	}

	mock.ExpectExec("INSERT INTO \"user\"").WithArgs(
		user.ID, user.Email, user.Name, user.Surname, user.PasswordHash, user.Version,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserCurrentVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db, logrus.New())

	userID := "test-user-id"

	mock.ExpectQuery("SELECT version FROM \"user\"").WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"version"}).AddRow(1),
	)

	version, err := repo.GetUserCurrentVersion(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, version)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db, logrus.New())

	email := "test@example.com"

	mock.ExpectQuery("SELECT user_id, email, name, surname, password_hash, version FROM \"user\"").WithArgs(email).WillReturnRows(
		sqlmock.NewRows([]string{"user_id", "email", "name", "surname", "password_hash", "version"}).
			AddRow(uuid.New().String(), email, "Test", "User", "hashedpassword", 1),
	)

	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, email, user.Email)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db, logrus.New())

	userID := uuid.New()
	expectedUserID := userID.String()

	mock.ExpectQuery("SELECT user_id, email, name, surname, password_hash, version FROM \"user\"").WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"user_id", "email", "name", "surname", "password_hash", "version"}).
			AddRow(expectedUserID, "test@example.com", "Test", "User", "hashedpassword", 1),
	)

	user, err := repo.GetUserByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, expectedUserID, user.ID.String())

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCheckUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewUserRepository(db, logrus.New())

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

	log := logrus.New()

	repo := repository.NewUserRepository(db, log)

	userID := "123"
	expectedQuery := `UPDATE \"user\" SET version = version \+ 1 WHERE user_id = \$1`

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

		assert.True(t, errors.Is(err, models.ErrUserNotFound))

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

	log := logrus.New()

	repo := repository.NewUserRepository(db, log)

	userID := "123"
	version := 5
	expectedQuery := `SELECT version FROM "user" WHERE user_id = \$1`

	t.Run("VersionMatches", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"version"}).AddRow(version)
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnRows(rows)

		result := repo.CheckUserVersion(context.Background(), userID, version)

		assert.True(t, result)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("VersionDoesNotMatch", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"version"}).AddRow(version + 1)
		mock.ExpectQuery(expectedQuery).
			WithArgs(userID).
			WillReturnRows(rows)

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
