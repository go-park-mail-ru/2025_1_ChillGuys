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
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

const (
	queryCreateUser        = `INSERT INTO bazaar.user (id, email, name, surname, password_hash, image_url) VALUES($1, $2, $3, $4, $5, $6);`
	queryCreateBasket      = `INSERT INTO bazaar.basket (id, user_id, total_price, total_price_discount) SELECT $1, $2, 0, 0`
	queryCreateUserVersion = `INSERT INTO bazaar.user_version (id, user_id, version, updated_at) VALUES($1, $2, $3, $4);`
)

func TestCreateUser_Success(t *testing.T) {
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
		PasswordHash: []byte("hashed_password"),
		ImageURL:     null.StringFrom("image.jpg"),
		UserVersion: models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.user").
		WithArgs(
			user.ID,
			user.Email,
			user.Name,
			user.Surname,
			user.PasswordHash,
			user.ImageURL,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO bazaar.user_version").
		WithArgs(
			user.UserVersion.ID,
			user.UserVersion.UserID,
			user.UserVersion.Version,
			user.UserVersion.UpdatedAt,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO bazaar.basket").
		WithArgs(sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	repo := auth2.NewAuthRepository(db)
	err = repo.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_TransactionError(t *testing.T) {
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
		PasswordHash: []byte("hashed_password"),
		ImageURL:     null.StringFrom("image.jpg"),
		UserVersion: models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

	repo := auth2.NewAuthRepository(db)
	err = repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "transaction error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_InsertUserError(t *testing.T) {
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
		PasswordHash: []byte("hashed_password"),
		ImageURL:     null.StringFrom("image.jpg"),
		UserVersion: models.UserVersionDB{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO bazaar.user").
		WithArgs(
			user.ID,
			user.Email,
			user.Name,
			user.Surname,
			user.PasswordHash,
			user.ImageURL,
		).
		WillReturnError(errors.New("insert user error"))
	mock.ExpectRollback()

	repo := auth2.NewAuthRepository(db)
	err = repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert user error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserCurrentVersion_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New().String()
	version := 1

	rows := sqlmock.NewRows([]string{"version"}).AddRow(version)

	mock.ExpectQuery("SELECT version FROM bazaar.user_version").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db)
	result, err := repo.GetUserCurrentVersion(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, version, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserCurrentVersion_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New().String()

	mock.ExpectQuery("SELECT version FROM bazaar.user_version").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	repo := auth2.NewAuthRepository(db)
	_, err = repo.GetUserCurrentVersion(context.Background(), userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	versionID := uuid.New()
	email := "test@example.com"
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "name", "surname", "password_hash", "image_url", "user_version_id", "version", "updated_at"}).
		AddRow(
			userID,
			email,
			"Test",
			"User",
			[]byte("hashed_password"),
			"image.jpg",
			versionID,
			1,
			updatedAt,
		)

	mock.ExpectQuery("SELECT").
		WithArgs(email).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db)
	user, err := repo.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, null.StringFrom("User"), user.Surname)
	assert.Equal(t, []byte("hashed_password"), user.PasswordHash)
	assert.Equal(t, null.StringFrom("image.jpg"), user.ImageURL)
	assert.Equal(t, versionID, user.UserVersion.ID)
	assert.Equal(t, 1, user.UserVersion.Version)
	assert.Equal(t, updatedAt, user.UserVersion.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "test@example.com"

	mock.ExpectQuery("SELECT").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	repo := auth2.NewAuthRepository(db)
	user, err := repo.GetUserByEmail(context.Background(), email)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	versionID := uuid.New()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "email", "name", "surname", "password_hash", "image_url", "phone_number", "user_version_id", "version", "updated_at"}).
		AddRow(
			userID,
			"test@example.com",
			"Test",
			"User",
			[]byte("hashed_password"),
			"image.jpg",
			"1234567890",
			versionID,
			1,
			updatedAt,
		)

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db)
	user, err := repo.GetUserByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, null.StringFrom("User"), user.Surname)
	assert.Equal(t, []byte("hashed_password"), user.PasswordHash)
	assert.Equal(t, null.StringFrom("image.jpg"), user.ImageURL)
	assert.Equal(t, null.StringFrom("1234567890"), user.PhoneNumber)
	assert.Equal(t, versionID, user.UserVersion.ID)
	assert.Equal(t, 1, user.UserVersion.Version)
	assert.Equal(t, updatedAt, user.UserVersion.UpdatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	repo := auth2.NewAuthRepository(db)
	user, err := repo.GetUserByID(context.Background(), userID)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIncrementUserVersion_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New().String()

	mock.ExpectExec("UPDATE bazaar.user_version").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := auth2.NewAuthRepository(db)
	err = repo.IncrementUserVersion(context.Background(), userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIncrementUserVersion_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New().String()

	mock.ExpectExec("UPDATE bazaar.user_version").
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := auth2.NewAuthRepository(db)
	err = repo.IncrementUserVersion(context.Background(), userID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckUserVersion_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New().String()
	version := 1

	rows := sqlmock.NewRows([]string{"version"}).AddRow(version)

	mock.ExpectQuery("SELECT version FROM bazaar.user_version").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db)
	result := repo.CheckUserVersion(context.Background(), userID, version)

	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckUserVersion_NotEqual(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New().String()
	currentVersion := 1
	checkVersion := 2

	rows := sqlmock.NewRows([]string{"version"}).AddRow(currentVersion)

	mock.ExpectQuery("SELECT version FROM bazaar.user_version").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db)
	result := repo.CheckUserVersion(context.Background(), userID, checkVersion)

	assert.False(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckUserExists_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "test@example.com"
	exists := true

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(exists)

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(email).
		WillReturnRows(rows)

	repo := auth2.NewAuthRepository(db)
	result, err := repo.CheckUserExists(context.Background(), email)

	assert.NoError(t, err)
	assert.True(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckUserExists_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "test@example.com"

	mock.ExpectQuery("SELECT EXISTS").
		WithArgs(email).
		WillReturnError(errors.New("database error"))

	repo := auth2.NewAuthRepository(db)
	_, err = repo.CheckUserExists(context.Background(), email)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
