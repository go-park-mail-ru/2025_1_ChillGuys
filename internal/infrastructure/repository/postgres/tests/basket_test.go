package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	basketRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/basket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasketRepository_Get(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := basketRepo.NewBasketRepository(db)

	t.Run("success with items", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		productID := uuid.New()
		now := time.Now()

		// First expectation - get basket ID
		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		// Second expectation - get basket items
		rows := sqlmock.NewRows([]string{
			"id", "basket_id", "product_id", "quantity", "updated_at",
			"name", "price", "preview_image_url", "discounted_price",
		}).AddRow(
			uuid.New(), basketID, productID, 2, now,
			"Test Product", 1000.0, "image.jpg", 800.0,
		)

		mock.ExpectQuery(`SELECT`).
			WithArgs(basketID).
			WillReturnRows(rows)

		items, err := repo.Get(context.Background(), userID)
		require.NoError(t, err)
		require.Len(t, items, 1)
	})

	t.Run("basket not found", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.Get(context.Background(), userID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("database error", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		_, err := repo.Get(context.Background(), userID)
		require.Error(t, err)
	})
}

func TestBasketRepository_Add(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := basketRepo.NewBasketRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		productID := uuid.New()
		now := time.Now()

		// First expectation - get basket ID
		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		// Second expectation - add product to basket
		// Use sqlmock.AnyArg() for the generated UUID
		mock.ExpectQuery(`INSERT INTO bazaar.basket_item`).
			WithArgs(sqlmock.AnyArg(), basketID, productID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "basket_id", "product_id", "quantity", "updated_at"}).
				AddRow(uuid.New(), basketID, productID, 1, now))

		_, err := repo.Add(context.Background(), userID, productID)
		require.NoError(t, err)
	})

	t.Run("basket not found", func(t *testing.T) {
		userID := uuid.New()
		productID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.Add(context.Background(), userID, productID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("product not found", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		productID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		mock.ExpectQuery(`INSERT INTO bazaar.basket_item`).
			WithArgs(sqlmock.AnyArg(), basketID, productID).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.Add(context.Background(), userID, productID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})
}

func TestBasketRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := basketRepo.NewBasketRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		productID := uuid.New()
		deletedID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		mock.ExpectQuery(`DELETE FROM bazaar.basket_item`).
			WithArgs(basketID, productID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(deletedID))

		err := repo.Delete(context.Background(), userID, productID)
		require.NoError(t, err)
	})

	t.Run("basket not found", func(t *testing.T) {
		userID := uuid.New()
		productID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		err := repo.Delete(context.Background(), userID, productID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})

	t.Run("product not in basket", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		productID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		mock.ExpectQuery(`DELETE FROM bazaar.basket_item`).
			WithArgs(basketID, productID).
			WillReturnError(sql.ErrNoRows)

		err := repo.Delete(context.Background(), userID, productID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})
}

func TestBasketRepository_UpdateQuantity(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := basketRepo.NewBasketRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		productID := uuid.New()
		itemID := uuid.New()
		now := time.Now()
		quantity := 2
		availableQuantity := uint(10)

		mock.ExpectQuery(`SELECT quantity FROM bazaar.product WHERE id = \$1`).
			WithArgs(productID).
			WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(availableQuantity))

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		mock.ExpectQuery(`UPDATE bazaar.basket_item`).
			WithArgs(quantity, basketID, productID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "basket_id", "product_id", "quantity", "updated_at"}).
				AddRow(itemID, basketID, productID, quantity, now))

		item, remaining, err := repo.UpdateQuantity(context.Background(), userID, productID, quantity)
		require.NoError(t, err)
		assert.Equal(t, itemID, item.ID)
		assert.Equal(t, quantity, item.Quantity)
		assert.Equal(t, int(availableQuantity)-quantity, remaining)
	})

	t.Run("insufficient quantity", func(t *testing.T) {
		userID := uuid.New()
		productID := uuid.New()
		quantity := 5
		availableQuantity := uint(3)

		mock.ExpectQuery(`SELECT quantity FROM bazaar.product WHERE id = \$1`).
			WithArgs(productID).
			WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(availableQuantity))

		_, _, err := repo.UpdateQuantity(context.Background(), userID, productID, quantity)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrBusinessLogic)
	})

	t.Run("product not found", func(t *testing.T) {
		userID := uuid.New()
		productID := uuid.New()

		mock.ExpectQuery(`SELECT quantity FROM bazaar.product WHERE id = \$1`).
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows)

		_, _, err := repo.UpdateQuantity(context.Background(), userID, productID, 1)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})
}

func TestBasketRepository_Clear(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := basketRepo.NewBasketRepository(db)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		mock.ExpectExec(`DELETE FROM bazaar.basket_item WHERE basket_id = \$1`).
			WithArgs(basketID).
			WillReturnResult(sqlmock.NewResult(0, 3))

		err := repo.Clear(context.Background(), userID)
		require.NoError(t, err)
	})

	t.Run("basket not found", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		err := repo.Clear(context.Background(), userID)
		require.Error(t, err)
		assert.ErrorIs(t, err, errs.ErrNotFound)
	})
}

func TestBasketRepository_GetProductsInBasket(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := basketRepo.NewBasketRepository(db)

	t.Run("success with discount", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()
		product1ID := uuid.New()
		now := time.Now()

		expectedItems := []*models.BasketItem{
			{
				ID:            uuid.New(),
				BasketID:      basketID,
				ProductID:     product1ID,
				Quantity:      2,
				UpdatedAt:     now,
				ProductName:   "Product 1",
				Price:         1000,
				ProductImage:  "image1.jpg",
				PriceDiscount: 800,
			},
		}

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		rows := sqlmock.NewRows([]string{
			"id", "basket_id", "product_id", "quantity", "updated_at",
			"name", "price", "preview_image_url", "discounted_price",
		}).
			AddRow(
				expectedItems[0].ID,
				expectedItems[0].BasketID,
				expectedItems[0].ProductID,
				expectedItems[0].Quantity,
				expectedItems[0].UpdatedAt,
				expectedItems[0].ProductName,
				expectedItems[0].Price,
				expectedItems[0].ProductImage,
				800.0,
			)

		mock.ExpectQuery(`SELECT`).
			WithArgs(basketID).
			WillReturnRows(rows)

		items, err := repo.Get(context.Background(), userID)
		require.NoError(t, err)
		require.Len(t, items, 1)
		assert.Equal(t, expectedItems, items)
	})

	t.Run("empty basket", func(t *testing.T) {
		userID := uuid.New()
		basketID := uuid.New()

		mock.ExpectQuery(`SELECT id FROM bazaar.basket WHERE user_id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(basketID))

		mock.ExpectQuery(`SELECT`).
			WithArgs(basketID).
			WillReturnRows(sqlmock.NewRows([]string{
				"id", "basket_id", "product_id", "quantity", "updated_at",
				"name", "price", "preview_image_url", "discounted_price",
			}))

		items, err := repo.Get(context.Background(), userID)
		require.NoError(t, err)
		assert.Empty(t, items)
	})
}