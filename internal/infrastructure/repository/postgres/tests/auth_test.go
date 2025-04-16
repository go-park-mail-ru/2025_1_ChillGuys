package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	auth2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	queryCreateUser        = `INSERT INTO bazaar.user (id, email, name, surname, password_hash, image_url) VALUES($1, $2, $3, $4, $5, $6);`
	queryCreateBasket      = `INSERT INTO bazaar.basket (id, user_id, total_price, total_price_discount) SELECT $1, $2, 0, 0`
	queryCreateUserVersion = `INSERT INTO bazaar.user_version (id, user_id, version, updated_at) VALUES($1, $2, $3, $4);`
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := models.UserDB{
		ID:           uuid.MustParse("2de435bd-5b89-4725-a03a-bab88141d607"),
		Email:        "test@example.com",
		Name:         "Test",
		Surname:      null.StringFrom("User"),
		PasswordHash: []byte("hashedpassword"),
		ImageURL:     null.StringFrom("http://example.com/image.jpg"),
		UserVersion: models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    uuid.MustParse("2de435bd-5b89-4725-a03a-bab88141d607"),
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec(queryCreateUser).
		WithArgs(
			user.ID,
			user.Email,
			user.Name,
			user.Surname,
			user.PasswordHash,
			user.ImageURL,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(queryCreateUserVersion).
		WithArgs(
			user.UserVersion.ID,
			user.UserVersion.UserID,
			user.UserVersion.Version,
			user.UserVersion.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(queryCreateBasket).
		WithArgs(sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := auth2.NewAuthRepository(db, logrus.New())
	err = repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_ErrorOnBeginTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	user := models.UserDB{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Name:         "Test",
		Surname:      null.StringFrom("User"),
		PasswordHash: []byte("hashedpassword"),
		ImageURL:     null.StringFrom("http://example.com/image.jpg"),
	}

	mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

	repo := auth2.NewAuthRepository(db, logrus.New())
	err = repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	expectedUser := models.UserDB{
		ID: userID,
		// ... other fields
		UserVersion: models.UserVersionDB{
			UserID: userID, // Make sure this matches the user's ID
			// ... other fields
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "email", "name", "surname", "password_hash", "image_url",
		"user_version_id", "version", "updated_at",
	}).AddRow(
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

	mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
		WithArgs(expectedUser.Email).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db, logrus.New())
	user, err := repo.GetUserByEmail(context.Background(), expectedUser.Email)

	assert.NoError(t, err)
	assert.Equal(t, &expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "notfound@example.com"

	mock.ExpectQuery(`SELECT.*FROM bazaar.user`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	repo := auth2.NewAuthRepository(db, logrus.New())
	user, err := repo.GetUserByEmail(context.Background(), email)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found") // or whatever error message your implementation uses
	assert.NoError(t, mock.ExpectationsWereMet())
}
