package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	productRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepository_GetAllProducts(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := productRepo.NewProductRepository(db)
	//
	//t.Run("success", func(t *testing.T) {
	//	product1ID := uuid.New()
	//	product2ID := uuid.New()
	//	sellerID := uuid.New()
	//	now := time.Now()
	//
	//	expectedProducts := []*models.Product{
	//		{
	//			ID:              product1ID,
	//			SellerID:        sellerID,
	//			Name:           "Product 1",
	//			PreviewImageURL: "image1.jpg",
	//			Description:    "Description 1",
	//			Status:         models.ProductApproved,
	//			Price:          1000,
	//			Quantity:       10,
	//			UpdatedAt:      now,
	//			Rating:         4,
	//			ReviewsCount:   20,
	//		},
	//		{
	//			ID:              product2ID,
	//			SellerID:        sellerID,
	//			Name:           "Product 2",
	//			PreviewImageURL: "image2.jpg",
	//			Description:    "Description 2",
	//			Status:         models.ProductApproved,
	//			Price:          2000,
	//			Quantity:       5,
	//			UpdatedAt:      now,
	//			Rating:         4,
	//			ReviewsCount:   15,
	//		},
	//	}
	//
	//	rows := sqlmock.NewRows([]string{
	//		"id", "seller_id", "name", "preview_image_url", "description",
	//		"status", "price", "quantity", "updated_at", "rating", "reviews_count",
	//	}).
	//		AddRow(
	//			expectedProducts[0].ID,
	//			expectedProducts[0].SellerID,
	//			expectedProducts[0].Name,
	//			expectedProducts[0].PreviewImageURL,
	//			expectedProducts[0].Description,
	//			expectedProducts[0].Status.String(),
	//			expectedProducts[0].Price,
	//			expectedProducts[0].Quantity,
	//			expectedProducts[0].UpdatedAt,
	//			expectedProducts[0].Rating,
	//			expectedProducts[0].ReviewsCount,
	//		).
	//		AddRow(
	//			expectedProducts[1].ID,
	//			expectedProducts[1].SellerID,
	//			expectedProducts[1].Name,
	//			expectedProducts[1].PreviewImageURL,
	//			expectedProducts[1].Description,
	//			expectedProducts[1].Status.String(),
	//			expectedProducts[1].Price,
	//			expectedProducts[1].Quantity,
	//			expectedProducts[1].UpdatedAt,
	//			expectedProducts[1].Rating,
	//			expectedProducts[1].ReviewsCount,
	//		)
	//
	//	mock.ExpectQuery(`
	//		SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description,
	//			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count
	//		FROM bazaar.product p
	//		WHERE p.status = 'approved'
	//	`).WillReturnRows(rows)
	//
	//	products, err := repo.GetAllProducts(context.Background())
	//	require.NoError(t, err)
	//	require.Len(t, products, 2)
	//	assert.Equal(t, expectedProducts, products)
	//})
	//
	//t.Run("empty result", func(t *testing.T) {
	//	rows := sqlmock.NewRows([]string{
	//		"id", "seller_id", "name", "preview_image_url", "description",
	//		"status", "price", "quantity", "updated_at", "rating", "reviews_count",
	//	})
	//
	//	mock.ExpectQuery(`
	//		SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description,
	//			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count
	//		FROM bazaar.product p
	//		WHERE p.status = 'approved'
	//	`).WillReturnRows(rows)
	//
	//	products, err := repo.GetAllProducts(context.Background())
	//	require.NoError(t, err)
	//	assert.Empty(t, products)
	//})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count 
			FROM bazaar.product p 
			WHERE p.status = 'approved'
		`).WillReturnError(errors.New("database error"))

		products, err := repo.GetAllProducts(context.Background())
		require.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestProductRepository_GetProductByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := productRepo.NewProductRepository(db)
	//
	//t.Run("success", func(t *testing.T) {
	//	productID := uuid.New()
	//	sellerID := uuid.New()
	//	now := time.Now()
	//
	//	expectedProduct := &models.Product{
	//		ID:              productID,
	//		SellerID:        sellerID,
	//		Name:           "Test Product",
	//		PreviewImageURL: "test.jpg",
	//		Description:    "Test Description",
	//		Status:         models.ProductApproved,
	//		Price:          1500,
	//		Quantity:       8,
	//		UpdatedAt:      now,
	//		Rating:         4,
	//		ReviewsCount:   10,
	//	}
	//
	//	row := sqlmock.NewRows([]string{
	//		"id", "seller_id", "name", "preview_image_url", "description",
	//		"status", "price", "quantity", "updated_at", "rating", "reviews_count",
	//	}).
	//		AddRow(
	//			expectedProduct.ID,
	//			expectedProduct.SellerID,
	//			expectedProduct.Name,
	//			expectedProduct.PreviewImageURL,
	//			expectedProduct.Description,
	//			expectedProduct.Status.String(),
	//			expectedProduct.Price,
	//			expectedProduct.Quantity,
	//			expectedProduct.UpdatedAt,
	//			expectedProduct.Rating,
	//			expectedProduct.ReviewsCount,
	//		)
	//
	//	mock.ExpectQuery(`
	//		SELECT id, seller_id, name, preview_image_url, description,
	//			status, price, quantity, updated_at, rating, reviews_count
	//		FROM bazaar.product WHERE id = \$1
	//	`).
	//		WithArgs(productID).
	//		WillReturnRows(row)
	//
	//	product, err := repo.GetProductByID(context.Background(), productID)
	//	require.NoError(t, err)
	//	assert.Equal(t, expectedProduct, product)
	//})
	//
	//t.Run("not found", func(t *testing.T) {
	//	productID := uuid.New()
	//
	//	mock.ExpectQuery(`
	//		SELECT id, seller_id, name, preview_image_url, description,
	//			status, price, quantity, updated_at, rating, reviews_count
	//		FROM bazaar.product WHERE id = \$1
	//	`).
	//		WithArgs(productID).
	//		WillReturnError(sql.ErrNoRows)
	//
	//	product, err := repo.GetProductByID(context.Background(), productID)
	//	require.Error(t, err)
	//	require.ErrorIs(t, err, errs.ErrNotFound)
	//	assert.Nil(t, product)
	//})

	t.Run("database error", func(t *testing.T) {
		productID := uuid.New()

		mock.ExpectQuery(`
			SELECT id, seller_id, name, preview_image_url, description, 
				status, price, quantity, updated_at, rating, reviews_count 
			FROM bazaar.product WHERE id = \$1
		`).
			WithArgs(productID).
			WillReturnError(errors.New("database error"))

		product, err := repo.GetProductByID(context.Background(), productID)
		require.Error(t, err)
		assert.Nil(t, product)
	})
}

func TestProductRepository_GetProductsByCategory(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := productRepo.NewProductRepository(db)

	t.Run("success", func(t *testing.T) {
		categoryID := uuid.New()
		product1ID := uuid.New()
		product2ID := uuid.New()
		sellerID := uuid.New()
		now := time.Now()

		expectedProducts := []*models.Product{
			{
				ID:              product1ID,
				SellerID:        sellerID,
				Name:            "Product 1",
				PreviewImageURL: "image1.jpg",
				Description:     "Description 1",
				Status:          models.ProductApproved,
				Price:           1000,
				Quantity:        10,
				UpdatedAt:       now,
				Rating:          4,
				ReviewsCount:    20,
			},
			{
				ID:              product2ID,
				SellerID:        sellerID,
				Name:            "Product 2",
				PreviewImageURL: "image2.jpg",
				Description:     "Description 2",
				Status:          models.ProductApproved,
				Price:           2000,
				Quantity:        5,
				UpdatedAt:       now,
				Rating:          4,
				ReviewsCount:    15,
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
		}).
			AddRow(
				expectedProducts[0].ID,
				expectedProducts[0].SellerID,
				expectedProducts[0].Name,
				expectedProducts[0].PreviewImageURL,
				expectedProducts[0].Description,
				expectedProducts[0].Status.String(),
				expectedProducts[0].Price,
				expectedProducts[0].Quantity,
				expectedProducts[0].UpdatedAt,
				expectedProducts[0].Rating,
				expectedProducts[0].ReviewsCount,
			).
			AddRow(
				expectedProducts[1].ID,
				expectedProducts[1].SellerID,
				expectedProducts[1].Name,
				expectedProducts[1].PreviewImageURL,
				expectedProducts[1].Description,
				expectedProducts[1].Status.String(),
				expectedProducts[1].Price,
				expectedProducts[1].Quantity,
				expectedProducts[1].UpdatedAt,
				expectedProducts[1].Rating,
				expectedProducts[1].ReviewsCount,
			)

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count 
			FROM bazaar.product p
			JOIN bazaar.product_category pc ON p.id = pc.product_id
			WHERE pc.category_id = \$1 AND p.status = 'approved'
		`).
			WithArgs(categoryID).
			WillReturnRows(rows)

		products, err := repo.GetProductsByCategory(context.Background(), categoryID)
		require.NoError(t, err)
		require.Len(t, products, 2)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("empty result", func(t *testing.T) {
		categoryID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
		})

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count 
			FROM bazaar.product p
			JOIN bazaar.product_category pc ON p.id = pc.product_id
			WHERE pc.category_id = \$1 AND p.status = 'approved'
		`).
			WithArgs(categoryID).
			WillReturnRows(rows)

		products, err := repo.GetProductsByCategory(context.Background(), categoryID)
		require.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("database error", func(t *testing.T) {
		categoryID := uuid.New()

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count 
			FROM bazaar.product p
			JOIN bazaar.product_category pc ON p.id = pc.product_id
			WHERE pc.category_id = \$1 AND p.status = 'approved'
		`).
			WithArgs(categoryID).
			WillReturnError(errors.New("database error"))

		products, err := repo.GetProductsByCategory(context.Background(), categoryID)
		require.Error(t, err)
		assert.Nil(t, products)
	})
}
