package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	review "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/review"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAdd_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIReviewRepository(ctrl)
	usecase := review.NewReviewUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := dto.AddReviewRequest{
		ProductID: uuid.New(),
		Rating:    5,
		Comment:   "Great product!",
	}

	mockRepo.EXPECT().
		AddReview(ctx, gomock.Any()).
		Do(func(_ context.Context, review models.ReviewDB) {
			assert.Equal(t, userID, review.UserID)
			assert.Equal(t, req.ProductID, review.ProductID)
			assert.Equal(t, req.Rating, review.Rating)
			assert.Equal(t, req.Comment, review.Comment)
		}).
		Return(nil)

	err := usecase.Add(ctx, req)

	assert.NoError(t, err)
}

func TestAdd_UserIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIReviewRepository(ctrl)
	usecase := review.NewReviewUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	req := dto.AddReviewRequest{
		ProductID: uuid.New(),
		Rating:    5,
		Comment:   "Great product!",
	}

	err := usecase.Add(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrNotFound))
}

func TestAdd_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIReviewRepository(ctrl)
	usecase := review.NewReviewUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	ctx = context.WithValue(ctx, domains.UserIDKey{}, "invalid-uuid")

	req := dto.AddReviewRequest{
		ProductID: uuid.New(),
		Rating:    5,
		Comment:   "Great product!",
	}

	err := usecase.Add(ctx, req)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidID))
}

func TestAdd_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIReviewRepository(ctrl)
	usecase := review.NewReviewUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))
	userID := uuid.New()
	ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())

	req := dto.AddReviewRequest{
		ProductID: uuid.New(),
		Rating:    5,
		Comment:   "Great product!",
	}

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		AddReview(ctx, gomock.Any()).
		Return(expectedError)

	err := usecase.Add(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIReviewRepository(ctrl)
	usecase := review.NewReviewUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	productID := uuid.New()
	offset := 0
	req := dto.GetReviewRequest{
		ProductID: productID,
		Offset:    offset,
	}

	expectedReviews := []*models.Review{
		{
			ID:      uuid.New(),
			Rating:  5,
			Comment: "Great product!",
			Name: "user1",
		},
		{
			ID:      uuid.New(),
			Rating:  4,
			Comment: "Good product",
			Name: "user2",
		},
	}

	mockRepo.EXPECT().
		GetReview(ctx, productID, offset).
		Return(expectedReviews, nil)

	reviews, err := usecase.Get(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedReviews, reviews)
}

func TestGet_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIReviewRepository(ctrl)
	usecase := review.NewReviewUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	productID := uuid.New()
	offset := 0
	req := dto.GetReviewRequest{
		ProductID: productID,
		Offset:    offset,
	}

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		GetReview(ctx, productID, offset).
		Return(nil, expectedError)

	_, err := usecase.Get(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}