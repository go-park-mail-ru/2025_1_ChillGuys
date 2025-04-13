package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	categoryRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/category"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepository_GetAllCategories(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := categoryRepo.NewCategoryRepository(db)

	t.Run("success", func(t *testing.T) {
		category1ID := uuid.New()
		category2ID := uuid.New()

		expectedCategories := []*models.Category{
			{
				ID:   category1ID,
				Name: "Category 1",
			},
			{
				ID:   category2ID,
				Name: "Category 2",
			},
		}

		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(category1ID, "Category 1").
			AddRow(category2ID, "Category 2")

		mock.ExpectQuery(`SELECT id, name FROM bazaar.category`).WillReturnRows(rows)

		categories, err := repo.GetAllCategories(context.Background())
		require.NoError(t, err)
		require.Len(t, categories, 2)
		assert.Equal(t, expectedCategories, categories)
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"})

		mock.ExpectQuery(`SELECT id, name FROM bazaar.category`).WillReturnRows(rows)

		categories, err := repo.GetAllCategories(context.Background())
		require.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, name FROM bazaar.category`).
			WillReturnError(errors.New("database error"))

		categories, err := repo.GetAllCategories(context.Background())
		require.Error(t, err)
		assert.Nil(t, categories)
	})

	t.Run("scan error", func(t *testing.T) {
		categoryID := uuid.New()

		// Return invalid data (missing name column)
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(categoryID)

		mock.ExpectQuery(`SELECT id, name FROM bazaar.category`).WillReturnRows(rows)

		categories, err := repo.GetAllCategories(context.Background())
		require.Error(t, err)
		assert.Nil(t, categories)
	})
}