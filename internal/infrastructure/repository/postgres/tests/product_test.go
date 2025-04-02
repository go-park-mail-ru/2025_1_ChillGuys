package tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := product.NewProductRepository(db, logrus.New())

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
		}).
			AddRow(
				uuid.New(), uuid.New(), "Product 1", "image1.jpg", "Description 1",
				"approved", 1000, 10, time.Now(), 4, 20,
			).
			AddRow(
				uuid.New(), uuid.New(), "Product 2", "image2.jpg", "Description 2",
				"approved", 2000, 5, time.Now(), 4, 15,
			)

		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product").
			WillReturnRows(rows)

		products, err := repo.GetAllProducts(context.Background())

		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, "Product 1", products[0].Name)
		assert.Equal(t, "Product 2", products[1].Name)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product").
			WillReturnError(errors.New("database error"))

		products, err := repo.GetAllProducts(context.Background())

		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestGetProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := product.NewProductRepository(db, logrus.New())
	testID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
		}).
			AddRow(
				testID, uuid.New(), "Test Product", "test.jpg", "Test Description",
				"approved", 1500, 8, time.Now(), 4, 10,
			)

		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product WHERE id = \\$1").
			WithArgs(testID).
			WillReturnRows(row)

		product, err := repo.GetProductByID(context.Background(), testID)

		assert.NoError(t, err)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, testID, product.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product WHERE id = \\$1").
			WithArgs(testID).
			WillReturnError(sql.ErrNoRows)

		product, err := repo.GetProductByID(context.Background(), testID)

		assert.Error(t, err)
		assert.Nil(t, product)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product WHERE id = \\$1").
			WithArgs(testID).
			WillReturnError(errors.New("database error"))

		product, err := repo.GetProductByID(context.Background(), testID)

		assert.Error(t, err)
		assert.Nil(t, product)
	})
}
