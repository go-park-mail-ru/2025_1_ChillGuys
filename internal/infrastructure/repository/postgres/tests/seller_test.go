package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	seller "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/seller"
)

func TestSellerRepository_AddProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := seller.NewSellerRepository(db)

	t.Run("Success", func(t *testing.T) {
		product := &models.Product{
			SellerID:  uuid.New(),
			Name:      "Test Product",
			Description: "Test Description",
			Price:     100.0,
			Quantity:  10,
		}
		categoryID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO bazaar.product").
			WithArgs(
				sqlmock.AnyArg(), // product.ID (будет сгенерирован)
				product.SellerID,
				product.Name,
				product.Description,
				models.ProductPending,
				product.Price,
				product.Quantity,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.product_subcategory").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), categoryID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		result, err := repo.AddProduct(context.Background(), product, categoryID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil, result.ID)
		assert.Equal(t, models.ProductPending, result.Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InsertProductError", func(t *testing.T) {
		product := &models.Product{
			SellerID:  uuid.New(),
			Name:      "Test Product",
			Description: "Test Description",
			Price:     100.0,
			Quantity:  10,
		}
		categoryID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO bazaar.product").
			WithArgs(
				sqlmock.AnyArg(),
				product.SellerID,
				product.Name,
				product.Description,
				models.ProductPending,
				product.Price,
				product.Quantity,
			).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		_, err := repo.AddProduct(context.Background(), product, categoryID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("InsertCategoryError", func(t *testing.T) {
		product := &models.Product{
			SellerID:  uuid.New(),
			Name:      "Test Product",
			Description: "Test Description",
			Price:     100.0,
			Quantity:  10,
		}
		categoryID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO bazaar.product").
			WithArgs(
				sqlmock.AnyArg(),
				product.SellerID,
				product.Name,
				product.Description,
				models.ProductPending,
				product.Price,
				product.Quantity,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO bazaar.product_subcategory").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), categoryID).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectRollback()

		_, err := repo.AddProduct(context.Background(), product, categoryID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_UploadProductImage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := seller.NewSellerRepository(db)

	t.Run("Success", func(t *testing.T) {
		productID := uuid.New()
		imageURL := "http://example.com/image.jpg"

		mock.ExpectExec("UPDATE bazaar.product").
			WithArgs(imageURL, productID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UploadProductImage(context.Background(), productID, imageURL)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error", func(t *testing.T) {
		productID := uuid.New()
		imageURL := "http://example.com/image.jpg"

		mock.ExpectExec("UPDATE bazaar.product").
			WithArgs(imageURL, productID).
			WillReturnError(sql.ErrConnDone)

		err := repo.UploadProductImage(context.Background(), productID, imageURL)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_GetSellerProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := seller.NewSellerRepository(db)

	t.Run("Success", func(t *testing.T) {
		sellerID := uuid.New()
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", 
			"description", "status", "price", "quantity", "rating", "reviews_count",
		}).
			AddRow(
				uuid.New(), sellerID, "Product 1", "image1.jpg", 
				"Description 1", models.ProductApproved, 100.0, 10, 4.5, 20,
			).
			AddRow(
				uuid.New(), sellerID, "Product 2", "image2.jpg", 
				"Description 2", models.ProductApproved, 200.0, 20, 4.0, 15,
			)

		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url,").
			WithArgs(sellerID, offset).
			WillReturnRows(rows)

		products, err := repo.GetSellerProducts(context.Background(), sellerID, offset)
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, "Product 1", products[0].Name)
		assert.Equal(t, "Product 2", products[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		sellerID := uuid.New()
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", 
			"description", "status", "price", "quantity", "rating", "reviews_count",
		})

		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url,").
			WithArgs(sellerID, offset).
			WillReturnRows(rows)

		products, err := repo.GetSellerProducts(context.Background(), sellerID, offset)
		assert.NoError(t, err)
		assert.Empty(t, products)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		sellerID := uuid.New()
		offset := 0

		mock.ExpectQuery("SELECT id, seller_id, name, preview_image_url,").
			WithArgs(sellerID, offset).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.GetSellerProducts(context.Background(), sellerID, offset)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSellerRepository_CheckProductBelongs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := seller.NewSellerRepository(db)

	t.Run("Belongs", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()

		rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(productID, sellerID).
			WillReturnRows(rows)

		belongs, err := repo.CheckProductBelongs(context.Background(), productID, sellerID)
		assert.NoError(t, err)
		assert.True(t, belongs)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotBelongs", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()

		rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(productID, sellerID).
			WillReturnRows(rows)

		belongs, err := repo.CheckProductBelongs(context.Background(), productID, sellerID)
		assert.NoError(t, err)
		assert.False(t, belongs)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("QueryError", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()

		mock.ExpectQuery("SELECT EXISTS").
			WithArgs(productID, sellerID).
			WillReturnError(sql.ErrConnDone)

		_, err := repo.CheckProductBelongs(context.Background(), productID, sellerID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}