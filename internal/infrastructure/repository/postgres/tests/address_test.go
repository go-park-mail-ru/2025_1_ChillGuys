package tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	address2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckAddressExists_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	addressID := uuid.New()
	address := models.AddressDB{
		AddressString: null.StringFrom("Test Address"),
		Coordinate:    null.StringFrom("10.0,20.0"),
	}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(addressID)

	mock.ExpectQuery("SELECT id FROM bazaar.address").
		WithArgs(address.AddressString, address.Coordinate).
		WillReturnRows(rows)

	repo := address2.NewAddressRepository(db)
	result, err := repo.CheckAddressExists(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, addressID, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckAddressExists_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	address := models.AddressDB{
		AddressString: null.StringFrom("Test Address"),
		Coordinate:    null.StringFrom("10.0,20.0"),
	}

	mock.ExpectQuery("SELECT id FROM bazaar.address").
		WithArgs(address.AddressString, address.Coordinate).
		WillReturnError(sql.ErrNoRows)

	repo := address2.NewAddressRepository(db)
	result, err := repo.CheckAddressExists(context.Background(), address)

	assert.NoError(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckAddressExists_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	address := models.AddressDB{
		AddressString: null.StringFrom("Test Address"),
		Coordinate:    null.StringFrom("10.0,20.0"),
	}

	mock.ExpectQuery("SELECT id FROM bazaar.address").
		WithArgs(address.AddressString, address.Coordinate).
		WillReturnError(errors.New("database error"))

	repo := address2.NewAddressRepository(db)
	result, err := repo.CheckAddressExists(context.Background(), address)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, result)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAddress_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	address := models.AddressDB{
		ID:            uuid.New(),
		Region:        null.StringFrom("Test Region"),
		City:          null.StringFrom("Test City"),
		AddressString: null.StringFrom("Test Address"),
		Coordinate:    null.StringFrom("10.0,20.0"),
	}

	mock.ExpectQuery("INSERT INTO bazaar.address").
		WithArgs(
			address.ID.String(),
			address.Region,
			address.City,
			address.AddressString,
			address.Coordinate,
		).
		WillReturnRows(sqlmock.NewRows([]string{}))

	repo := address2.NewAddressRepository(db)
	err = repo.CreateAddress(context.Background(), address)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAddress_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	address := models.AddressDB{
		ID:            uuid.New(),
		Region:        null.StringFrom("Test Region"),
		City:          null.StringFrom("Test City"),
		AddressString: null.StringFrom("Test Address"),
		Coordinate:    null.StringFrom("10.0,20.0"),
	}

	mock.ExpectQuery("INSERT INTO bazaar.address").
		WithArgs(
			address.ID.String(),
			address.Region,
			address.City,
			address.AddressString,
			address.Coordinate,
		).
		WillReturnError(errors.New("database error"))

	repo := address2.NewAddressRepository(db)
	err = repo.CreateAddress(context.Background(), address)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserAddress_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userAddress := models.UserAddress{
		ID:        uuid.New(),
		Label:     null.StringFrom("Home"),
		UserID:    uuid.New(),
		AddressID: uuid.New(),
	}

	mock.ExpectExec("INSERT INTO bazaar.user_address").
		WithArgs(
			userAddress.ID.String(),
			userAddress.Label,
			userAddress.UserID.String(),
			userAddress.AddressID.String(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := address2.NewAddressRepository(db)
	err = repo.CreateUserAddress(context.Background(), userAddress)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserAddress_ExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userAddress := models.UserAddress{
		ID:        uuid.New(),
		Label:     null.StringFrom("Home"),
		UserID:    uuid.New(),
		AddressID: uuid.New(),
	}

	mock.ExpectExec("INSERT INTO bazaar.user_address").
		WithArgs(
			userAddress.ID.String(),
			userAddress.Label,
			userAddress.UserID.String(),
			userAddress.AddressID.String(),
		).
		WillReturnError(errors.New("database error"))

	repo := address2.NewAddressRepository(db)
	err = repo.CreateUserAddress(context.Background(), userAddress)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserAddress_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()
	addressID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "label", "region", "city", "address_string", "coordinate"}).
		AddRow(
			addressID,
			"Home",
			"Test Region",
			"Test City",
			"Test Address",
			"10.0,20.0",
		)

	mock.ExpectQuery("SELECT a.id, ua.label, a.region, a.city, a.address_string, a.coordinate").
		WithArgs(userID.String()).
		WillReturnRows(rows)

	repo := address2.NewAddressRepository(db)
	addresses, err := repo.GetUserAddress(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, *addresses, 1)
	assert.Equal(t, addressID, (*addresses)[0].ID)
	assert.Equal(t, null.StringFrom("Home"), (*addresses)[0].Label)
	assert.Equal(t, null.StringFrom("Test Region"), (*addresses)[0].Region)
	assert.Equal(t, null.StringFrom("Test City"), (*addresses)[0].City)
	assert.Equal(t, null.StringFrom("Test Address"), (*addresses)[0].AddressString)
	assert.Equal(t, null.StringFrom("10.0,20.0"), (*addresses)[0].Coordinate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserAddress_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	mock.ExpectQuery("SELECT a.id, ua.label, a.region, a.city, a.address_string, a.coordinate").
		WithArgs(userID.String()).
		WillReturnError(errors.New("database error"))

	repo := address2.NewAddressRepository(db)
	addresses, err := repo.GetUserAddress(context.Background(), userID)

	assert.Nil(t, addresses)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserAddress_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "label", "region", "city", "address_string", "coordinate"}).
		AddRow("invalid_uuid", "Home", "Test Region", "Test City", "Test Address", "10.0,20.0")

	mock.ExpectQuery("SELECT a.id, ua.label, a.region, a.city, a.address_string, a.coordinate").
		WithArgs(userID.String()).
		WillReturnRows(rows)

	repo := address2.NewAddressRepository(db)
	addresses, err := repo.GetUserAddress(context.Background(), userID)

	assert.Nil(t, addresses)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPickupPoints_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	addressID := uuid.New()

	rows := sqlmock.NewRows([]string{"id", "region", "city", "address_string", "coordinate"}).
		AddRow(
			addressID,
			"Test Region",
			"Test City",
			"Test Address",
			"10.0,20.0",
		)

	mock.ExpectQuery("SELECT a.id, a.region, a.city, a.address_string, a.coordinate").
		WillReturnRows(rows)

	repo := address2.NewAddressRepository(db)
	points, err := repo.GetAllPickupPoints(context.Background())

	assert.NoError(t, err)
	assert.Len(t, *points, 1)
	assert.Equal(t, addressID, (*points)[0].ID)
	assert.Equal(t, null.StringFrom("Test Region"), (*points)[0].Region)
	assert.Equal(t, null.StringFrom("Test City"), (*points)[0].City)
	assert.Equal(t, null.StringFrom("Test Address"), (*points)[0].AddressString)
	assert.Equal(t, null.StringFrom("10.0,20.0"), (*points)[0].Coordinate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPickupPoints_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT a.id, a.region, a.city, a.address_string, a.coordinate").
		WillReturnError(errors.New("database error"))

	repo := address2.NewAddressRepository(db)
	points, err := repo.GetAllPickupPoints(context.Background())

	assert.Nil(t, points)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPickupPoints_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "region", "city", "address_string", "coordinate"}).
		AddRow("invalid_uuid", "Test Region", "Test City", "Test Address", "10.0,20.0")

	mock.ExpectQuery("SELECT a.id, a.region, a.city, a.address_string, a.coordinate").
		WillReturnRows(rows)

	repo := address2.NewAddressRepository(db)
	points, err := repo.GetAllPickupPoints(context.Background())

	assert.Nil(t, points)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
