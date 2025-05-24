package tests

import (
	"context"
	"database/sql"
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

	t.Run("success", func(t *testing.T) {
		product1ID := uuid.New()
		product2ID := uuid.New()
		sellerID := uuid.New()
		now := time.Now()
		offset := 0

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
				PriceDiscount:   800,
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
				PriceDiscount:   0,
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
			"discounted_price",
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
				expectedProducts[0].PriceDiscount,
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
				nil,
			)

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price
			FROM bazaar.product p 
			LEFT JOIN bazaar.discount d ON p.id = d.product_id
			WHERE p.status = 'approved'
			ORDER BY p.id
			LIMIT 20 OFFSET \$1
		`).WithArgs(offset).WillReturnRows(rows)

		products, err := repo.GetAllProducts(context.Background(), offset)
		require.NoError(t, err)
		require.Len(t, products, 2)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("database error", func(t *testing.T) {
		offset := 0
		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price
			FROM bazaar.product p 
			LEFT JOIN bazaar.discount d ON p.id = d.product_id
			WHERE p.status = 'approved'
			ORDER BY p.id
			LIMIT 20 OFFSET \$1
		`).WithArgs(offset).WillReturnError(errors.New("database error"))

		products, err := repo.GetAllProducts(context.Background(), offset)
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

	t.Run("success with seller", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()
		now := time.Now()

		expectedProduct := &models.Product{
			ID:              productID,
			SellerID:        sellerID,
			Name:            "Test Product",
			PreviewImageURL: "test.jpg",
			Description:     "Test Description",
			Status:          models.ProductApproved,
			Price:           1500,
			Quantity:        8,
			UpdatedAt:       now,
			Rating:          4,
			ReviewsCount:    10,
			PriceDiscount:   1200,
			Seller: &models.Seller{
				ID:          sellerID,
				Title:       "Test Seller",
				Description: "Test Description",
			},
		}

		row := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
			"discounted_price", "id", "title", "description",
		}).
			AddRow(
				expectedProduct.ID,
				expectedProduct.SellerID,
				expectedProduct.Name,
				expectedProduct.PreviewImageURL,
				expectedProduct.Description,
				expectedProduct.Status.String(),
				expectedProduct.Price,
				expectedProduct.Quantity,
				expectedProduct.UpdatedAt,
				expectedProduct.Rating,
				expectedProduct.ReviewsCount,
				expectedProduct.PriceDiscount,
				expectedProduct.Seller.ID,
				expectedProduct.Seller.Title,
				expectedProduct.Seller.Description,
			)

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price, s.id, s.title, s.description
			FROM bazaar.product p
			LEFT JOIN bazaar.discount d ON p.id = d.product_id
			LEFT JOIN bazaar.seller s ON s.user_id = p.seller_id
			WHERE p.id = \$1
		`).
			WithArgs(productID).
			WillReturnRows(row)

		product, err := repo.GetProductByID(context.Background(), productID)
		require.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("not found", func(t *testing.T) {
		productID := uuid.New()

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price, s.id, s.title, s.description
			FROM bazaar.product p
			LEFT JOIN bazaar.discount d ON p.id = d.product_id
			LEFT JOIN bazaar.seller s ON s.user_id = p.seller_id
			WHERE p.id = \$1
		`).
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows)

		product, err := repo.GetProductByID(context.Background(), productID)
		require.Error(t, err)
		require.ErrorContains(t, err, "not found") // More flexible error check
		assert.Nil(t, product)
	})

	t.Run("database error", func(t *testing.T) {
		productID := uuid.New()

		mock.ExpectQuery(`
			SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price, s.id, s.title, s.description
			FROM bazaar.product p
			LEFT JOIN bazaar.discount d ON p.id = d.product_id
			LEFT JOIN bazaar.seller s ON s.user_id = p.seller_id
			WHERE p.id = \$1
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

	t.Run("success with price_asc sort", func(t *testing.T) {
		categoryID := uuid.New()
		product1ID := uuid.New()
		product2ID := uuid.New()
		sellerID := uuid.New()
		now := time.Now()
		offset := 0
		minPrice := 1000.0
		maxPrice := 2000.0
		minRating := float32(4.0)

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
			SELECT 
				p.id, 
				p.seller_id, 
				p.name, 
				p.preview_image_url, 
				p.description, 
				p.status, 
				p.price, 
				p.quantity, 
				p.updated_at, 
				p.rating, 
				p.reviews_count
			FROM 
				bazaar.product p
			JOIN 
				bazaar.product_subcategory pc ON p.id = pc.product_id
			WHERE 
				pc.subcategory_id = \$1 
				AND p.status = 'approved'
				AND \(\$3 = 0 OR p.price > \$3\)
				AND \(\$4 = 0 OR p.price < \$4\)
				AND \(\$5 = 0::FLOAT OR p.rating > \$5::FLOAT\)
			ORDER BY p.price ASC
			LIMIT 20 OFFSET \$2
		`).
			WithArgs(categoryID, offset, minPrice, maxPrice, minRating).
			WillReturnRows(rows)

		products, err := repo.GetProductsByCategory(
			context.Background(),
			categoryID,
			offset,
			minPrice,
			maxPrice,
			minRating,
			models.SortByPriceAsc,
		)
		require.NoError(t, err)
		require.Len(t, products, 2)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("empty result", func(t *testing.T) {
		categoryID := uuid.New()
		offset := 0
		minPrice := 0.0
		maxPrice := 0.0
		minRating := float32(0.0)

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
		})

		mock.ExpectQuery(`
			SELECT 
				p.id, 
				p.seller_id, 
				p.name, 
				p.preview_image_url, 
				p.description, 
				p.status, 
				p.price, 
				p.quantity, 
				p.updated_at, 
				p.rating, 
				p.reviews_count
			FROM 
				bazaar.product p
			JOIN 
				bazaar.product_subcategory pc ON p.id = pc.product_id
			WHERE 
				pc.subcategory_id = \$1 
				AND p.status = 'approved'
				AND \(\$3 = 0 OR p.price > \$3\)
				AND \(\$4 = 0 OR p.price < \$4\)
				AND \(\$5 = 0::FLOAT OR p.rating > \$5::FLOAT\)
			ORDER BY p.updated_at DESC
			LIMIT 20 OFFSET \$2
		`).
			WithArgs(categoryID, offset, minPrice, maxPrice, minRating).
			WillReturnRows(rows)

		products, err := repo.GetProductsByCategory(
			context.Background(),
			categoryID,
			offset,
			minPrice,
			maxPrice,
			minRating,
			models.SortByDefault,
		)
		require.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("database error", func(t *testing.T) {
		categoryID := uuid.New()
		offset := 0
		minPrice := 0.0
		maxPrice := 0.0
		minRating := float32(0.0)

		mock.ExpectQuery(`
			SELECT 
				p.id, 
				p.seller_id, 
				p.name, 
				p.preview_image_url, 
				p.description, 
				p.status, 
				p.price, 
				p.quantity, 
				p.updated_at, 
				p.rating, 
				p.reviews_count
			FROM 
				bazaar.product p
			JOIN 
				bazaar.product_subcategory pc ON p.id = pc.product_id
			WHERE 
				pc.subcategory_id = \$1 
				AND p.status = 'approved'
				AND \(\$3 = 0 OR p.price > \$3\)
				AND \(\$4 = 0 OR p.price < \$4\)
				AND \(\$5 = 0::FLOAT OR p.rating > \$5::FLOAT\)
			ORDER BY p.updated_at DESC
			LIMIT 20 OFFSET \$2
		`).
			WithArgs(categoryID, offset, minPrice, maxPrice, minRating).
			WillReturnError(errors.New("database error"))

		products, err := repo.GetProductsByCategory(
			context.Background(),
			categoryID,
			offset,
			minPrice,
			maxPrice,
			minRating,
			models.SortByDefault,
		)
		require.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestProductRepository_AddProduct(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := productRepo.NewProductRepository(db)

	t.Run("success", func(t *testing.T) {
		product := &models.Product{
			SellerID:        uuid.New(),
			Name:            "New Product",
			PreviewImageURL: "new.jpg",
			Description:     "New Description",
			Price:           2000,
			PriceDiscount:   1800,
			Quantity:        10,
			Rating:          0,
			ReviewsCount:    0,
		}
		categoryID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`
			INSERT INTO bazaar.product \(
				id, seller_id, name, preview_image_url, 
				description, status, price, quantity, rating, reviews_count
			\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$10, \$6, \$7, \$8, \$9\)
			RETURNING id
		`).
			WithArgs(
				sqlmock.AnyArg(),
				product.SellerID,
				product.Name,
				product.PreviewImageURL,
				product.Description,
				product.Price,
				product.Quantity,
				product.Rating,
				product.ReviewsCount,
				models.ProductApproved, // This matches the $10 parameter
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`
			INSERT INTO bazaar.discount \(
				id, product_id, discounted_price, start_date, end_date
			\) VALUES \(\$1, \$2, \$3, now\(\), now\(\) \+ interval '30 days'\)
		`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), product.PriceDiscount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`
			INSERT INTO bazaar.product_subcategory \(id, product_id, subcategory_id\)
			VALUES \(\$1, \$2, \$3\)
		`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), categoryID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		result, err := repo.AddProduct(context.Background(), product, categoryID)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, result.ID)
		assert.Equal(t, models.ProductApproved, result.Status)
	})

	t.Run("transaction error", func(t *testing.T) {
		product := &models.Product{
			SellerID:        uuid.New(),
			Name:            "New Product",
			PreviewImageURL: "new.jpg",
			Description:     "New Description",
			Price:           2000,
			PriceDiscount:   1800,
			Quantity:        10,
			Rating:          0,
			ReviewsCount:    0,
		}
		categoryID := uuid.New()

		mock.ExpectBegin().WillReturnError(errors.New("transaction error"))

		result, err := repo.AddProduct(context.Background(), product, categoryID)
		require.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("insert product error", func(t *testing.T) {
		product := &models.Product{
			SellerID:        uuid.New(),
			Name:            "New Product",
			PreviewImageURL: "new.jpg",
			Description:     "New Description",
			Price:           2000,
			PriceDiscount:   1800,
			Quantity:        10,
			Rating:          0,
			ReviewsCount:    0,
		}
		categoryID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`
			INSERT INTO bazaar.product \(
				id, seller_id, name, preview_image_url, 
				description, status, price, quantity, rating, reviews_count
			\) VALUES \(\$1, \$2, \$3, \$4, \$5, 'approved', \$6, \$7, \$8, \$9\)
			RETURNING id
		`).
			WithArgs(
				sqlmock.AnyArg(),
				product.SellerID,
				product.Name,
				product.PreviewImageURL,
				product.Description,
				product.Price,
				product.Quantity,
				product.Rating,
				product.ReviewsCount,
			).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		result, err := repo.AddProduct(context.Background(), product, categoryID)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}
