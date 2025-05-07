package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	search "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/search"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchRepository_GetCategoryByName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := search.NewSearchRepository(db)

	t.Run("success", func(t *testing.T) {
		categoryID := uuid.New()
		categoryName := "Test Category"

		expectedCategory := &models.Category{
			ID:   categoryID,
			Name: categoryName,
		}

		row := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(expectedCategory.ID, expectedCategory.Name)

		mock.ExpectQuery("SELECT").
			WithArgs(categoryName).
			WillReturnRows(row)

		category, err := repo.GetCategoryByName(context.Background(), categoryName)
		require.NoError(t, err)
		assert.Equal(t, expectedCategory, category)
	})

	t.Run("not found", func(t *testing.T) {
		categoryName := "Non-existent Category"

		mock.ExpectQuery("SELECT").
			WithArgs(categoryName).
			WillReturnError(sql.ErrNoRows)

		category, err := repo.GetCategoryByName(context.Background(), categoryName)
		require.NoError(t, err)
		assert.Nil(t, category)
	})

	t.Run("database error", func(t *testing.T) {
		categoryName := "Test Category"

		mock.ExpectQuery("SELECT").
			WithArgs(categoryName).
			WillReturnError(errors.New("database error"))

		category, err := repo.GetCategoryByName(context.Background(), categoryName)
		require.Error(t, err)
		assert.Nil(t, category)
	})
}

func TestSearchRepository_GetProductsByNameWithFilterAndSort(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := search.NewSearchRepository(db)

	t.Run("success with all filters and price_asc sort", func(t *testing.T) {
		product1ID := uuid.New()
		product2ID := uuid.New()
		sellerID := uuid.New()
		categoryID := uuid.New()
		now := time.Now()
		searchTerm := "test"
		offset := 0
		minPrice := 1000.0
		maxPrice := 2000.0
		minRating := float32(4.0)

		expectedProducts := []*models.Product{
			{
				ID:              product1ID,
				SellerID:        sellerID,
				Name:            "Test Product 1",
				PreviewImageURL: "image1.jpg",
				Description:     "Description 1",
				Status:          models.ProductApproved,
				Price:           1000,
				Quantity:        10,
				UpdatedAt:       now,
				Rating:          4.5,
				ReviewsCount:    20,
				PriceDiscount:   800,
			},
			{
				ID:              product2ID,
				SellerID:        sellerID,
				Name:            "Test Product 2",
				PreviewImageURL: "image2.jpg",
				Description:     "Description 2",
				Status:          models.ProductApproved,
				Price:           1500,
				Quantity:        5,
				UpdatedAt:       now,
				Rating:          4.0,
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

		// Используем регулярное выражение для игнорирования пробелов и переносов строк
		expectedQuery := `SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
			d.discounted_price
		FROM bazaar.product p
		JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
		WHERE p.status = 'approved'
		AND LOWER\(p.name\) LIKE LOWER\(\$1\)
		AND \(\$2 = '' OR ps.subcategory_id = \$2::uuid\)
		AND \(\$3 = 0 OR p.price >= \$3\)
		AND \(\$4 = 0 OR p.price <= \$4\)
		AND \(\$5 = 0::FLOAT OR p.rating >= \$5::FLOAT\)
		ORDER BY p.price ASC
		LIMIT 20 OFFSET \$6`

		mock.ExpectQuery(expectedQuery).
			WithArgs(
				"%"+searchTerm+"%",
				categoryID.String(),
				minPrice,
				maxPrice,
				minRating,
				offset,
			).
			WillReturnRows(rows)

		products, err := repo.GetProductsByNameWithFilterAndSort(
			context.Background(),
			searchTerm,
			null.StringFrom(categoryID.String()),
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

	t.Run("success with minimal filters and default sort", func(t *testing.T) {
		productID := uuid.New()
		sellerID := uuid.New()
		now := time.Now()
		searchTerm := "test"
		offset := 0

		expectedProduct := &models.Product{
			ID:              productID,
			SellerID:        sellerID,
			Name:            "Test Product",
			PreviewImageURL: "image.jpg",
			Description:     "Description",
			Status:          models.ProductApproved,
			Price:           1000,
			Quantity:        10,
			UpdatedAt:       now,
			Rating:          4.5,
			ReviewsCount:    20,
			PriceDiscount:   800,
		}

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
			"discounted_price",
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
			)

		expectedQuery := `SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
			d.discounted_price
		FROM bazaar.product p
		JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
		WHERE p.status = 'approved'
		AND LOWER\(p.name\) LIKE LOWER\(\$1\)
		AND \(\$2 = '' OR ps.subcategory_id = \$2::uuid\)
		AND \(\$3 = 0 OR p.price >= \$3\)
		AND \(\$4 = 0 OR p.price <= \$4\)
		AND \(\$5 = 0::FLOAT OR p.rating >= \$5::FLOAT\)
		ORDER BY p.updated_at DESC
		LIMIT 20 OFFSET \$6`

		mock.ExpectQuery(expectedQuery).
			WithArgs(
				"%"+searchTerm+"%",
				"",
				0.0,
				0.0,
				float32(0.0),
				offset,
			).
			WillReturnRows(rows)

		products, err := repo.GetProductsByNameWithFilterAndSort(
			context.Background(),
			searchTerm,
			null.String{},
			offset,
			0.0,
			0.0,
			0.0,
			models.SortByDefault,
		)
		require.NoError(t, err)
		require.Len(t, products, 1)
		assert.Equal(t, expectedProduct, products[0])
	})

	t.Run("empty result", func(t *testing.T) {
		searchTerm := "non-existent"
		offset := 0

		rows := sqlmock.NewRows([]string{
			"id", "seller_id", "name", "preview_image_url", "description",
			"status", "price", "quantity", "updated_at", "rating", "reviews_count",
			"discounted_price",
		})

		expectedQuery := `SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
			d.discounted_price
		FROM bazaar.product p
		JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
		WHERE p.status = 'approved'
		AND LOWER\(p.name\) LIKE LOWER\(\$1\)
		AND \(\$2 = '' OR ps.subcategory_id = \$2::uuid\)
		AND \(\$3 = 0 OR p.price >= \$3\)
		AND \(\$4 = 0 OR p.price <= \$4\)
		AND \(\$5 = 0::FLOAT OR p.rating >= \$5::FLOAT\)
		ORDER BY p.updated_at DESC
		LIMIT 20 OFFSET \$6`

		mock.ExpectQuery(expectedQuery).
			WithArgs(
				"%"+searchTerm+"%",
				"",
				0.0,
				0.0,
				float32(0.0),
				offset,
			).
			WillReturnRows(rows)

		products, err := repo.GetProductsByNameWithFilterAndSort(
			context.Background(),
			searchTerm,
			null.String{},
			offset,
			0.0,
			0.0,
			0.0,
			models.SortByDefault,
		)
		require.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("database error", func(t *testing.T) {
		searchTerm := "test"
		offset := 0

		expectedQuery := `SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
			d.discounted_price
		FROM bazaar.product p
		JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
		WHERE p.status = 'approved'
		AND LOWER\(p.name\) LIKE LOWER\(\$1\)
		AND \(\$2 = '' OR ps.subcategory_id = \$2::uuid\)
		AND \(\$3 = 0 OR p.price >= \$3\)
		AND \(\$4 = 0 OR p.price <= \$4\)
		AND \(\$5 = 0::FLOAT OR p.rating >= \$5::FLOAT\)
		ORDER BY p.updated_at DESC
		LIMIT 20 OFFSET \$6`

		mock.ExpectQuery(expectedQuery).
			WithArgs(
				"%"+searchTerm+"%",
				"",
				0.0,
				0.0,
				float32(0.0),
				offset,
			).
			WillReturnError(errors.New("database error"))

		products, err := repo.GetProductsByNameWithFilterAndSort(
			context.Background(),
			searchTerm,
			null.String{},
			offset,
			0.0,
			0.0,
			0.0,
			models.SortByDefault,
		)
		require.Error(t, err)
		assert.Nil(t, products)
	})
}