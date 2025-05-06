package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	suggestions "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/suggestions"
)

func TestSuggestionsRepository_GetAllCategoriesName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := suggestions.NewSuggestionsRepository(db)

	t.Run("success", func(t *testing.T) {
		expectedCategories := []*models.CategorySuggestion{
			{Name: "Category 1"},
			{Name: "Category 2"},
		}

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Category 1").
			AddRow("Category 2")

		mock.ExpectQuery(`SELECT name FROM bazaar.subcategory`).
			WillReturnRows(rows)

		categories, err := repo.GetAllCategoriesName(context.Background())
		require.NoError(t, err)
		require.Len(t, categories, 2)
		assert.Equal(t, expectedCategories, categories)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(`SELECT name FROM bazaar.subcategory`).
			WillReturnRows(rows)

		categories, err := repo.GetAllCategoriesName(context.Background())
		require.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT name FROM bazaar.subcategory`).
			WillReturnError(errors.New("database error"))

		categories, err := repo.GetAllCategoriesName(context.Background())
		require.Error(t, err)
		assert.Nil(t, categories)
	})

	t.Run("scan error", func(t *testing.T) {
		// Return more columns than expected
		rows := sqlmock.NewRows([]string{"name", "extra_column"}).
			AddRow("Category 1", "extra")

		mock.ExpectQuery(`SELECT name FROM bazaar.subcategory`).
			WillReturnRows(rows)

		categories, err := repo.GetAllCategoriesName(context.Background())
		require.Error(t, err)
		assert.Nil(t, categories)
	})
}

func TestSuggestionsRepository_GetAllProductsName(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := suggestions.NewSuggestionsRepository(db)

	t.Run("success", func(t *testing.T) {
		expectedProducts := []*models.ProductSuggestion{
			{Name: "Product 1"},
			{Name: "Product 2"},
		}

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Product 1").
			AddRow("Product 2")

		mock.ExpectQuery(`SELECT name FROM bazaar.product WHERE status = 'approved'`).
			WillReturnRows(rows)

		products, err := repo.GetAllProductsName(context.Background())
		require.NoError(t, err)
		require.Len(t, products, 2)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(`SELECT name FROM bazaar.product WHERE status = 'approved'`).
			WillReturnRows(rows)

		products, err := repo.GetAllProductsName(context.Background())
		require.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT name FROM bazaar.product WHERE status = 'approved'`).
			WillReturnError(errors.New("database error"))

		products, err := repo.GetAllProductsName(context.Background())
		require.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestSuggestionsRepository_GetProductsNameByCategory(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := suggestions.NewSuggestionsRepository(db)
	categoryID := "123e4567-e89b-12d3-a456-426614174000"

	t.Run("success", func(t *testing.T) {
		expectedProducts := []*models.ProductSuggestion{
			{Name: "Product 1"},
			{Name: "Product 2"},
		}

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Product 1").
			AddRow("Product 2")

		mock.ExpectQuery(`
			SELECT DISTINCT p.name
			FROM bazaar.product p
			JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
			JOIN bazaar.subcategory s ON s.id = ps.subcategory_id
			WHERE s.id = \$1 AND p.status = 'approved'`).
			WithArgs(categoryID).
			WillReturnRows(rows)

		products, err := repo.GetProductsNameByCategory(context.Background(), categoryID)
		require.NoError(t, err)
		require.Len(t, products, 2)
		assert.Equal(t, expectedProducts, products)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery(`
			SELECT DISTINCT p.name
			FROM bazaar.product p
			JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
			JOIN bazaar.subcategory s ON s.id = ps.subcategory_id
			WHERE s.id = \$1 AND p.status = 'approved'`).
			WithArgs(categoryID).
			WillReturnRows(rows)

		products, err := repo.GetProductsNameByCategory(context.Background(), categoryID)
		require.NoError(t, err)
		assert.Empty(t, products)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`
			SELECT DISTINCT p.name
			FROM bazaar.product p
			JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
			JOIN bazaar.subcategory s ON s.id = ps.subcategory_id
			WHERE s.id = \$1 AND p.status = 'approved'`).
			WithArgs(categoryID).
			WillReturnError(errors.New("database error"))

		products, err := repo.GetProductsNameByCategory(context.Background(), categoryID)
		require.Error(t, err)
		assert.Nil(t, products)
	})

	t.Run("invalid category id format", func(t *testing.T) {
		invalidCategoryID := "invalid-uuid"

		// Expect no query execution for invalid input
		products, err := repo.GetProductsNameByCategory(context.Background(), invalidCategoryID)
		require.Error(t, err)
		assert.Nil(t, products)
	})
}