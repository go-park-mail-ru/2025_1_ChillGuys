package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/promo"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePromoCode_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := models.PromoCode{
		ID:        uuid.New(),
		Code:      "SUMMER20",
		Percent:   20,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec("INSERT INTO bazaar.promo_code").
		WithArgs(p.ID, p.Code, p.Percent, p.StartDate, p.EndDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := promo.NewPromoRepository(db)
	err = repo.Create(context.Background(), p)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePromoCode_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	p := models.PromoCode{
		ID:        uuid.New(),
		Code:      "SUMMER20",
		Percent:   20,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	mock.ExpectExec("INSERT INTO bazaar.promo_code").
		WithArgs(p.ID, p.Code, p.Percent, p.StartDate, p.EndDate).
		WillReturnError(errors.New("database error"))

	repo := promo.NewPromoRepository(db)
	err = repo.Create(context.Background(), p)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPromoCodes_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	expected := []*models.PromoCode{
		{
			ID:        uuid.New(),
			Code:      "SUMMER20",
			Percent:   20,
			StartDate: now,
			EndDate:   now.Add(24 * time.Hour),
		},
		{
			ID:        uuid.New(),
			Code:      "WINTER30",
			Percent:   30,
			StartDate: now.Add(-48 * time.Hour),
			EndDate:   now.Add(-24 * time.Hour),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "code", "percent", "start_date", "end_date"}).
		AddRow(expected[0].ID, expected[0].Code, expected[0].Percent, expected[0].StartDate, expected[0].EndDate).
		AddRow(expected[1].ID, expected[1].Code, expected[1].Percent, expected[1].StartDate, expected[1].EndDate)

	mock.ExpectQuery("SELECT id, code, percent, start_date, end_date FROM bazaar.promo_code").
		WillReturnRows(rows)

	repo := promo.NewPromoRepository(db)
	result, err := repo.GetAll(context.Background(), 0)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPromoCodes_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "code", "percent", "start_date", "end_date"})

	mock.ExpectQuery("SELECT id, code, percent, start_date, end_date FROM bazaar.promo_code").
		WillReturnRows(rows)

	repo := promo.NewPromoRepository(db)
	result, err := repo.GetAll(context.Background(), 0)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPromoCodes_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, code, percent, start_date, end_date FROM bazaar.promo_code").
		WillReturnError(errors.New("database error"))

	repo := promo.NewPromoRepository(db)
	_, err = repo.GetAll(context.Background(), 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckPromoCode_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	expected := &models.PromoCode{
		ID:        uuid.New(),
		Code:      "SUMMER20",
		Percent:   20,
		StartDate: now,
		EndDate:   now.Add(24 * time.Hour),
	}

	rows := sqlmock.NewRows([]string{"id", "code", "percent", "start_date", "end_date"}).
		AddRow(expected.ID, expected.Code, expected.Percent, expected.StartDate, expected.EndDate)

	mock.ExpectQuery("SELECT id, code, percent, start_date, end_date FROM bazaar.promo_code WHERE code = ?").
		WithArgs("SUMMER20").
		WillReturnRows(rows)

	repo := promo.NewPromoRepository(db)
	result, err := repo.CheckPromoCode(context.Background(), "SUMMER20")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckPromoCode_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, code, percent, start_date, end_date FROM bazaar.promo_code WHERE code = ?").
		WithArgs("ERROR").
		WillReturnError(errors.New("database error"))

	repo := promo.NewPromoRepository(db)
	_, err = repo.CheckPromoCode(context.Background(), "ERROR")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
