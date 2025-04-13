package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

//go:generate mockgen -source=product.go -destination=../../infrastructure/repository/postgres/mocks/product_repository_mock.go -package=mocks IProductRepository
type IProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductsByCategory(ctx context.Context, id uuid.UUID) ([]*models.Product, error)
}

type ProductUsecase struct {
	repo IProductRepository
}

func NewProductUsecase(repo IProductRepository) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	const op = "ProductUsecase.GetAllProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	products, err := u.repo.GetAllProducts(ctx)
	if err != nil {
		logger.WithError(err).Error("get products from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	const op = "ProductUsecase.GetProductByID"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_id", id)

	product, err := u.repo.GetProductByID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("get product by ID from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return product, nil
}

func (u *ProductUsecase) GetProductsByCategory(ctx context.Context, id uuid.UUID) ([]*models.Product, error) {
	const op = "ProductUsecase.GetProductsByCategory"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("category_id", id)

	products, err := u.repo.GetProductsByCategory(ctx, id)
	if err != nil {
		logger.WithError(err).Error("get products by category from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

// GetProductsByIDs возвращает список продуктов по их UUID
func (u *ProductUsecase) GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Product, error) {
	const op = "ProductUsecase.GetProductsByIDs"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_ids", ids)

	if len(ids) == 0 {
		return []*models.Product{}, nil
	}

	var products []*models.Product
	for _, id := range ids {
		product, err := u.repo.GetProductByID(ctx, id)
		if err != nil {
			logger.WithError(err).WithField("product_id", id).Warn("failed to get product by ID")
			continue
		}
		products = append(products, product)
	}

	return products, nil
}
