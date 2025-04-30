package seller

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

//go:generate mockgen -source=seller.go -destination=../../infrastructure/repository/postgres/mocks/seller_repository_mock.go -package=mocks ISellerRepository
type ISellerRepository interface {
	AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error)
	UploadProductImage(ctx context.Context, productID uuid.UUID, imageURL string) error
	GetSellerProducts(ctx context.Context, sellerID uuid.UUID, offset int) ([]*models.Product, error)
	CheckProductBelongs(ctx context.Context, productID, sellerID uuid.UUID) (bool, error)
}

type SellerUsecase struct {
	repo ISellerRepository
}

func NewSellerUsecase(repo ISellerRepository) *SellerUsecase {
	return &SellerUsecase{
		repo: repo,
	}
}

func (u *SellerUsecase) AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error) {
	const op = "SellerUsecase.AddProduct"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Валидация данных
	if product.Name == "" {
		logger.Error("empty product name")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrEmptyProductName)
	}
	if product.Price <= 0 {
		logger.Error("invalid product price")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidProductPrice)
	}

	// Добавляем продукт
	newProduct, err := u.repo.AddProduct(ctx, product, categoryID)
	if err != nil {
		logger.WithError(err).Error("add product to repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return newProduct, nil
}

func (u *SellerUsecase) UploadProductImage(ctx context.Context, productID uuid.UUID, imageURL string) error {
	const op = "SellerUsecase.UploadProductImage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	err := u.repo.UploadProductImage(ctx, productID, imageURL)
	if err != nil {
		logger.WithError(err).Error("upload product image")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *SellerUsecase) GetSellerProducts(ctx context.Context, sellerID uuid.UUID, offset int) ([]*models.Product, error) {
	const op = "SellerUsecase.GetSellerProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	products, err := u.repo.GetSellerProducts(ctx, sellerID, offset)
	if err != nil {
		logger.WithError(err).Error("get seller products from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (u *SellerUsecase) CheckProductBelongs(ctx context.Context, productID, sellerID uuid.UUID) (bool, error) {
	const op = "SellerUsecase.CheckProductBelongs"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	belongs, err := u.repo.CheckProductBelongs(ctx, productID, sellerID)
	if err != nil {
		logger.WithError(err).Error("check product belongs")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return belongs, nil
}