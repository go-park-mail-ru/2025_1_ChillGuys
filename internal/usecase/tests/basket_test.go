package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/basket"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, domains.UserIDKey{}, userID.String())
}

func setupTestBasket(t *testing.T) (*mocks.MockIBasketRepository, *basket.BasketUsecase) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockIBasketRepository(ctrl)
	uc := basket.NewBasketUsecase(mockRepo)
	return mockRepo, uc
}

func TestBasketUsecase_Get(t *testing.T) {
	userID := uuid.New()
	ctx := ContextWithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		expectedItems := []*models.BasketItem{
			{
				ID:        uuid.New(),
				BasketID:  uuid.New(),
				ProductID: uuid.New(),
				Quantity:  2,
			},
		}

		mockRepo.EXPECT().
			Get(gomock.Any(), userID).
			Return(expectedItems, nil)

		items, err := uc.Get(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedItems, items)
	})

	t.Run("no user in context", func(t *testing.T) {
		_, uc := setupTestBasket(t)
		_, err := uc.Get(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		mockRepo.EXPECT().
			Get(gomock.Any(), userID).
			Return(nil, errors.New("db error"))

		_, err := uc.Get(ctx)
		assert.Error(t, err)
	})
}

func TestBasketUsecase_Add(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	ctx := ContextWithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		expectedItem := &models.BasketItem{
			ID:        uuid.New(),
			BasketID:  uuid.New(),
			ProductID: productID,
			Quantity:  1,
		}

		mockRepo.EXPECT().
			Add(gomock.Any(), userID, productID).
			Return(expectedItem, nil)

		item, err := uc.Add(ctx, productID)
		assert.NoError(t, err)
		assert.Equal(t, expectedItem, item)
	})

	t.Run("invalid product id", func(t *testing.T) {
		_, uc := setupTestBasket(t)
		_, err := uc.Add(ctx, uuid.Nil)
		assert.ErrorIs(t, err, errs.ErrInvalidID)
	})

	// t.Run("no user in context", func(t *testing.T) {
	// 	_, uc := setupTestBasket(t)
	// 	_, err := uc.Add(context.Background(), productID)
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), errs.ErrNotFound)
	// })

	t.Run("repository error", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		mockRepo.EXPECT().
			Add(gomock.Any(), userID, productID).
			Return(nil, errors.New("db error"))

		_, err := uc.Add(ctx, productID)
		assert.Error(t, err)
	})
}

func TestBasketUsecase_Delete(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()
	ctx := ContextWithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		mockRepo.EXPECT().
			Delete(gomock.Any(), userID, productID).
			Return(nil)

		err := uc.Delete(ctx, productID)
		assert.NoError(t, err)
	})

	t.Run("no user in context", func(t *testing.T) {
		_, uc := setupTestBasket(t)
		err := uc.Delete(context.Background(), productID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		mockRepo.EXPECT().
			Delete(gomock.Any(), userID, productID).
			Return(errors.New("db error"))

		err := uc.Delete(ctx, productID)
		assert.Error(t, err)
	})
}

//func TestBasketUsecase_UpdateQuantity(t *testing.T) {
//	userID := uuid.New()
//	productID := uuid.New()
//	quantity := 2
//	ctx := ContextWithUserID(context.Background(), userID)

//t.Run("success", func(t *testing.T) {
//	mockRepo, uc := setupTestBasket(t)
//	expectedItem := &models.BasketItem{
//		ID:        uuid.New(),
//		BasketID:  uuid.New(),
//		ProductID: productID,
//		Quantity:  quantity,
//	}
//	remaining := 5
//
//	mockRepo.EXPECT().
//		UpdateQuantity(gomock.Any(), userID, productID, quantity).
//		Return(expectedItem, remaining, nil)
//
//	item, rem, err := uc.UpdateQuantity(ctx, productID, quantity)
//	assert.NoError(t, err)
//	assert.Equal(t, expectedItem, item)
//	assert.Equal(t, remaining, rem)
//})

//t.Run("invalid quantity", func(t *testing.T) {
//	_, uc := setupTestBasket(t)
//	_, _, err := uc.UpdateQuantity(ctx, productID, 0)
//	assert.Error(t, err)
//})
//
//t.Run("invalid product id", func(t *testing.T) {
//	_, uc := setupTestBasket(t)
//	_, _, err := uc.UpdateQuantity(ctx, uuid.Nil, quantity)
//	assert.ErrorIs(t, err, errs.ErrInvalidID)
//})
//
//t.Run("no user in context", func(t *testing.T) {
//	_, uc := setupTestBasket(t)
//	_, _, err := uc.UpdateQuantity(context.Background(), productID, quantity)
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "user not found")
//})

//	t.Run("repository error", func(t *testing.T) {
//		mockRepo, uc := setupTestBasket(t)
//		mockRepo.EXPECT().
//			UpdateQuantity(gomock.Any(), userID, productID, quantity).
//			Return(nil, -1, errors.New("db error"))
//
//		_, _, err := uc.UpdateQuantity(ctx, productID, quantity)
//		assert.Error(t, err)
//	})
//}

func TestBasketUsecase_Clear(t *testing.T) {
	userID := uuid.New()
	ctx := ContextWithUserID(context.Background(), userID)

	t.Run("success", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		mockRepo.EXPECT().
			Clear(gomock.Any(), userID).
			Return(nil)

		err := uc.Clear(ctx)
		assert.NoError(t, err)
	})

	t.Run("no user in context", func(t *testing.T) {
		_, uc := setupTestBasket(t)
		err := uc.Clear(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo, uc := setupTestBasket(t)
		mockRepo.EXPECT().
			Clear(gomock.Any(), userID).
			Return(errors.New("db error"))

		err := uc.Clear(ctx)
		assert.Error(t, err)
	})
}
