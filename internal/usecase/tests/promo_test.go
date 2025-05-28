package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/promo"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePromo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	req := dto.CreatePromoRequest{
		Code:      "SUMMER20",
		Percent:   20,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, p models.PromoCode) error {
			assert.Equal(t, req.Code, p.Code)
			assert.Equal(t, req.Percent, p.Percent)
			assert.Equal(t, req.StartDate, p.StartDate)
			assert.Equal(t, req.EndDate, p.EndDate)
			return nil
		})

	result, err := uc.CreatePromo(ctx, req)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, req.Code, result.Code)
	assert.Equal(t, req.Percent, result.Percent)
	assert.Equal(t, req.StartDate, result.StartDate)
	assert.Equal(t, req.EndDate, result.EndDate)
}

func TestCreatePromo_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	req := dto.CreatePromoRequest{
		Code:      "SUMMER20",
		Percent:   20,
		StartDate: time.Now(),
		EndDate:   time.Now().Add(24 * time.Hour),
	}

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	_, err := uc.CreatePromo(ctx, req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestGetAllPromos_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	now := time.Now()
	promos := []*models.PromoCode{
		{
			ID:        uuid.New(),
			Code:      "SUMMER20",
			Percent:   20,
			StartDate: now,
			EndDate:   now.Add(24 * time.Hour),
		},
		{
			ID:        uuid.New(),
			Code:      "WINTER30",
			Percent:   30,
			StartDate: now.Add(-48 * time.Hour),
			EndDate:   now.Add(-24 * time.Hour),
		},
	}

	mockRepo.EXPECT().
		GetAll(ctx, 0).
		Return(promos, nil)

	result, err := uc.GetAllPromos(ctx, 0)

	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Promos, 2)
	assert.Equal(t, promos[0].Code, result.Promos[0].Code)
	assert.Equal(t, promos[1].Code, result.Promos[1].Code)
}

func TestGetAllPromos_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	mockRepo.EXPECT().
		GetAll(ctx, 0).
		Return(nil, nil)

	_, err := uc.GetAllPromos(ctx, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetAllPromos_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	mockRepo.EXPECT().
		GetAll(ctx, 0).
		Return(nil, errors.New("database error"))

	_, err := uc.GetAllPromos(ctx, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestCheckPromoCode_Valid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	code := "SUMMER20"
	now := time.Now()
	promo := &models.PromoCode{
		ID:        uuid.New(),
		Code:      code,
		Percent:   20,
		StartDate: now.Add(-1 * time.Hour),
		EndDate:   now.Add(1 * time.Hour),
	}

	mockRepo.EXPECT().
		CheckPromoCode(ctx, code).
		Return(promo, nil)

	result, err := uc.CheckPromoCode(ctx, code)

	require.NoError(t, err)
	assert.True(t, result.IsValid)
	assert.Equal(t, promo.Percent, result.Percent)
}

func TestCheckPromoCode_Expired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	code := "SUMMER20"
	now := time.Now()
	promo := &models.PromoCode{
		ID:        uuid.New(),
		Code:      code,
		Percent:   20,
		StartDate: now.Add(-48 * time.Hour),
		EndDate:   now.Add(-24 * time.Hour),
	}

	mockRepo.EXPECT().
		CheckPromoCode(ctx, code).
		Return(promo, nil)

	result, err := uc.CheckPromoCode(ctx, code)

	require.NoError(t, err)
	assert.False(t, result.IsValid)
	assert.Equal(t, 0, result.Percent)
}

func TestCheckPromoCode_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIPromoRepository(ctrl)
	uc := promo.NewPromoUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	code := "INVALID"

	mockRepo.EXPECT().
		CheckPromoCode(ctx, code).
		Return(nil, errors.New("not found"))

	result, err := uc.CheckPromoCode(ctx, code)

	require.NoError(t, err)
	assert.False(t, result.IsValid)
	assert.Equal(t, 0, result.Percent)
}
