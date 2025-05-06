package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	categoryRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
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

func TestCategoryRepository_GetAllSubcategories(t *testing.T) {
    t.Parallel()

    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    repo := categoryRepo.NewCategoryRepository(db)

    t.Run("success", func(t *testing.T) {
        categoryID := uuid.New()
        subcategory1ID := uuid.New()
        subcategory2ID := uuid.New()

        expectedSubcategories := []*models.Category{
            {
                ID:   subcategory1ID,
                Name: "Subcategory 1",
            },
            {
                ID:   subcategory2ID,
                Name: "Subcategory 2",
            },
        }

        rows := sqlmock.NewRows([]string{"id", "name"}).
            AddRow(subcategory1ID, "Subcategory 1").
            AddRow(subcategory2ID, "Subcategory 2")

        mock.ExpectQuery(`SELECT id, name FROM bazaar.subcategory WHERE category_id = \$1`).
            WithArgs(categoryID).
            WillReturnRows(rows)

        subcategories, err := repo.GetAllSubcategories(context.Background(), categoryID)
        require.NoError(t, err)
        require.Len(t, subcategories, 2)
        assert.Equal(t, expectedSubcategories, subcategories)
    })

    t.Run("empty result", func(t *testing.T) {
        categoryID := uuid.New()

        rows := sqlmock.NewRows([]string{"id", "name"})

        mock.ExpectQuery(`SELECT id, name FROM bazaar.subcategory WHERE category_id = \$1`).
            WithArgs(categoryID).
            WillReturnRows(rows)

        subcategories, err := repo.GetAllSubcategories(context.Background(), categoryID)
        require.NoError(t, err)
        assert.Empty(t, subcategories)
    })

    t.Run("database error", func(t *testing.T) {
        categoryID := uuid.New()

        mock.ExpectQuery(`SELECT id, name FROM bazaar.subcategory WHERE category_id = \$1`).
            WithArgs(categoryID).
            WillReturnError(errors.New("database error"))

        subcategories, err := repo.GetAllSubcategories(context.Background(), categoryID)
        require.Error(t, err)
        assert.Nil(t, subcategories)
    })

    t.Run("scan error", func(t *testing.T) {
        categoryID := uuid.New()
        subcategoryID := uuid.New()

        // Return invalid data (missing name column)
        rows := sqlmock.NewRows([]string{"id"}).
            AddRow(subcategoryID)

        mock.ExpectQuery(`SELECT id, name FROM bazaar.subcategory WHERE category_id = \$1`).
            WithArgs(categoryID).
            WillReturnRows(rows)

        subcategories, err := repo.GetAllSubcategories(context.Background(), categoryID)
        require.Error(t, err)
        assert.Nil(t, subcategories)
    })
}

func TestCategoryRepository_GetNameSubcategory(t *testing.T) {
    t.Parallel()

    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    repo := categoryRepo.NewCategoryRepository(db)

    t.Run("success", func(t *testing.T) {
        subcategoryID := uuid.New()
        expectedName := "Test Subcategory"

        row := sqlmock.NewRows([]string{"name"}).
            AddRow(expectedName)

        mock.ExpectQuery(`SELECT name FROM bazaar.subcategory WHERE id = \$1`).
            WithArgs(subcategoryID).
            WillReturnRows(row)

        name, err := repo.GetNameSubcategory(context.Background(), subcategoryID)
        require.NoError(t, err)
        assert.Equal(t, expectedName, name)
    })

    t.Run("not found", func(t *testing.T) {
        subcategoryID := uuid.New()

        mock.ExpectQuery(`SELECT name FROM bazaar.subcategory WHERE id = \$1`).
            WithArgs(subcategoryID).
            WillReturnError(sql.ErrNoRows)

        name, err := repo.GetNameSubcategory(context.Background(), subcategoryID)
        require.Error(t, err)
        assert.Equal(t, "", name)
        assert.True(t, errors.Is(err, errs.ErrNotFound))
    })

    t.Run("database error", func(t *testing.T) {
        subcategoryID := uuid.New()

        mock.ExpectQuery(`SELECT name FROM bazaar.subcategory WHERE id = \$1`).
            WithArgs(subcategoryID).
            WillReturnError(errors.New("database error"))

        name, err := repo.GetNameSubcategory(context.Background(), subcategoryID)
        require.Error(t, err)
        assert.Equal(t, "", name)
    })
}