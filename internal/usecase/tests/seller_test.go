package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	seller "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/seller"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAddProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	testProduct := &models.Product{
		Name:  "Test Product",
		Price: 100,
	}
	categoryID := uuid.New()

	expectedProduct := &models.Product{
		ID:    uuid.New(),
		Name:  "Test Product",
		Price: 100,
	}

	mockRepo.EXPECT().
		AddProduct(ctx, testProduct, categoryID).
		Return(expectedProduct, nil)

	result, err := usecase.AddProduct(ctx, testProduct, categoryID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)
}

func TestAddProduct_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	testProduct := &models.Product{
		Name:  "",
		Price: 100,
	}
	categoryID := uuid.New()

	_, err := usecase.AddProduct(ctx, testProduct, categoryID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrEmptyProductName))
}

func TestAddProduct_InvalidPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	testProduct := &models.Product{
		Name:  "Test Product",
		Price: -100,
	}
	categoryID := uuid.New()

	_, err := usecase.AddProduct(ctx, testProduct, categoryID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrInvalidProductPrice))
}

func TestAddProduct_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	testProduct := &models.Product{
		Name:  "Test Product",
		Price: 100,
	}
	categoryID := uuid.New()

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		AddProduct(ctx, testProduct, categoryID).
		Return(nil, expectedError)

	_, err := usecase.AddProduct(ctx, testProduct, categoryID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestUploadProductImage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	productID := uuid.New()
	imageURL := "http://example.com/image.jpg"

	mockRepo.EXPECT().
		UploadProductImage(ctx, productID, imageURL).
		Return(nil)

	err := usecase.UploadProductImage(ctx, productID, imageURL)

	assert.NoError(t, err)
}

func TestUploadProductImage_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	productID := uuid.New()
	imageURL := "http://example.com/image.jpg"

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		UploadProductImage(ctx, productID, imageURL).
		Return(expectedError)

	err := usecase.UploadProductImage(ctx, productID, imageURL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestGetSellerProducts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	sellerID := uuid.New()
	offset := 0

	expectedProducts := []*models.Product{
		{
			ID:    uuid.New(),
			Name:  "Product 1",
			Price: 100,
		},
		{
			ID:    uuid.New(),
			Name:  "Product 2",
			Price: 200,
		},
	}

	mockRepo.EXPECT().
		GetSellerProducts(ctx, sellerID, offset).
		Return(expectedProducts, nil)

	result, err := usecase.GetSellerProducts(ctx, sellerID, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, result)
}

func TestGetSellerProducts_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	sellerID := uuid.New()
	offset := 0

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		GetSellerProducts(ctx, sellerID, offset).
		Return(nil, expectedError)

	_, err := usecase.GetSellerProducts(ctx, sellerID, offset)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestCheckProductBelongs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	productID := uuid.New()
	sellerID := uuid.New()

	expectedResult := true

	mockRepo.EXPECT().
		CheckProductBelongs(ctx, productID, sellerID).
		Return(expectedResult, nil)

	result, err := usecase.CheckProductBelongs(ctx, productID, sellerID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestCheckProductBelongs_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISellerRepository(ctrl)
	usecase := seller.NewSellerUsecase(mockRepo)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	productID := uuid.New()
	sellerID := uuid.New()

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().
		CheckProductBelongs(ctx, productID, sellerID).
		Return(false, expectedError)

	_, err := usecase.CheckProductBelongs(ctx, productID, sellerID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
}