package tests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	admin "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/admin"
)

func TestAdminRepository_GetPendingProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := admin.NewAdminRepository(db)

	t.Run("Success", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()
		now := time.Now()

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description", 
			"status", "price", "quantity", "updated_at", "rating", 
			"reviews_count", "discounted_price", "id", "title", "description",
		}).
			AddRow(
				productID, sellerID, "Test Product", "image.jpg", "Description",
				"pending", 100.0, 10, now, 4.5, 20, 90.0,
				sellerID, "Seller Name", "Seller Description",
			)

		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnRows(rows)

		products, err := repo.GetPendingProducts(context.Background(), 0)
		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, productID, products[0].ID)
		assert.Equal(t, models.ProductPending, products[0].Status)
		assert.Equal(t, 90.0, products[0].PriceDiscount)
		assert.NotNil(t, products[0].Seller)
		assert.Equal(t, sellerID, products[0].Seller.ID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description", 
			"status", "price", "quantity", "updated_at", "rating", 
			"reviews_count", "discounted_price", "id", "title", "description",
		})

		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnRows(rows)

		products, err := repo.GetPendingProducts(context.Background(), 0)
		assert.NoError(t, err)
		assert.Empty(t, products)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetPendingProducts(context.Background(), 0)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description", 
			"status", "price", "quantity", "updated_at", "rating", 
			"reviews_count", "discounted_price", "id", "title", "description",
		}).
			AddRow("invalid-uuid", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnRows(rows)

		_, err := repo.GetPendingProducts(context.Background(), 0)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_UpdateProductStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := admin.NewAdminRepository(db)
	productID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE").
			WithArgs("approved", productID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateProductStatus(context.Background(), productID, models.ProductApproved)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No rows affected", func(t *testing.T) {
		mock.ExpectExec("UPDATE").
			WithArgs("approved", productID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateProductStatus(context.Background(), productID, models.ProductApproved)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec error", func(t *testing.T) {
		mock.ExpectExec("UPDATE").
			WithArgs("approved", productID).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateProductStatus(context.Background(), productID, models.ProductApproved)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_GetPendingUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := admin.NewAdminRepository(db)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		sellerID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "image_url", 
			"phone_number", "role", "id", "title", "description",
		}).
			AddRow(
				userID, "test@example.com", "John", "Doe", "image.jpg",
				"+1234567890", "pending", sellerID, "Seller Name", "Seller Desc",
			)

		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnRows(rows)

		users, err := repo.GetPendingUsers(context.Background(), 0)
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, userID, users[0].ID)
		assert.Equal(t, models.RolePending, users[0].Role)
		assert.NotNil(t, users[0].Seller)
		assert.Equal(t, sellerID, users[0].Seller.ID)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "email", "name", "surname", "image_url", 
			"phone_number", "role", "id", "title", "description",
		})

		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnRows(rows)

		users, err := repo.GetPendingUsers(context.Background(), 0)
		assert.NoError(t, err)
		assert.Empty(t, users)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT").
			WithArgs(0).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetPendingUsers(context.Background(), 0)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAdminRepository_UpdateUserRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := admin.NewAdminRepository(db)
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("UPDATE").
			WithArgs("seller", userID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateUserRole(context.Background(), userID, models.RoleSeller)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No rows affected", func(t *testing.T) {
		mock.ExpectExec("UPDATE").
			WithArgs("seller", userID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUserRole(context.Background(), userID, models.RoleSeller)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec error", func(t *testing.T) {
		mock.ExpectExec("UPDATE").
			WithArgs("seller", userID).
			WillReturnError(sql.ErrConnDone)

		err := repo.UpdateUserRole(context.Background(), userID, models.RoleSeller)
		assert.Error(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}