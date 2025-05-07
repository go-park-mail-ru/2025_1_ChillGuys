package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
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
		Role:         models.RoleBuyer,
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
			user.Role,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO bazaar.basket").
		WithArgs(sqlmock.AnyArg(), user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	repo := auth.NewAuthRepository(db)
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
		Role:         models.RoleBuyer,
	}

	mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

	repo := auth.NewAuthRepository(db)
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
		Role:         models.RoleBuyer,
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
			user.Role,
		).
		WillReturnError(errors.New("insert user error"))
	mock.ExpectRollback()

	repo := auth.NewAuthRepository(db)
	err = repo.CreateUser(context.Background(), user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert user error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	email := "test@example.com"

	rows := sqlmock.NewRows([]string{"id", "email", "name", "surname", "password_hash", "image_url", "role"}).
		AddRow(
			userID,
			email,
			"Test",
			"User",
			[]byte("hashed_password"),
			"image.jpg",
			models.RoleBuyer,
		)

	mock.ExpectQuery("SELECT id, email, name, surname, password_hash, image_url, role FROM bazaar.user WHERE email =").
		WithArgs(email).
		WillReturnRows(rows)

	repo := auth.NewAuthRepository(db)
	user, err := repo.GetUserByEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, null.StringFrom("User"), user.Surname)
	assert.Equal(t, []byte("hashed_password"), user.PasswordHash)
	assert.Equal(t, null.StringFrom("image.jpg"), user.ImageURL)
	assert.Equal(t, models.RoleBuyer, user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	email := "test@example.com"

	mock.ExpectQuery("SELECT id, email, name, surname, password_hash, image_url, role FROM bazaar.user WHERE email =").
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	repo := auth.NewAuthRepository(db)
	user, err := repo.GetUserByEmail(context.Background(), email)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidCredentials))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "email", "name", "surname", "password_hash", "image_url", "phone_number", "role"}).
		AddRow(
			userID,
			"test@example.com",
			"Test",
			"User",
			[]byte("hashed_password"),
			"image.jpg",
			"1234567890",
			models.RoleBuyer,
		)

	mock.ExpectQuery("SELECT id, email, name, surname, password_hash, image_url, phone_number, role FROM bazaar.user WHERE id =").
		WithArgs(userID).
		WillReturnRows(rows)

	repo := auth.NewAuthRepository(db)
	user, err := repo.GetUserByID(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, null.StringFrom("User"), user.Surname)
	assert.Equal(t, []byte("hashed_password"), user.PasswordHash)
	assert.Equal(t, null.StringFrom("image.jpg"), user.ImageURL)
	assert.Equal(t, null.StringFrom("1234567890"), user.PhoneNumber)
	assert.Equal(t, models.RoleBuyer, user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT id, email, name, surname, password_hash, image_url, phone_number, role FROM bazaar.user WHERE id =").
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	repo := auth.NewAuthRepository(db)
	user, err := repo.GetUserByID(context.Background(), userID)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidCredentials))
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

	repo := auth.NewAuthRepository(db)
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

	repo := auth.NewAuthRepository(db)
	_, err = repo.CheckUserExists(context.Background(), email)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}