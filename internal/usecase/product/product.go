package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

//go:generate mockgen -source=product.go -destination=../../infrastructure/repository/postgres/mocks/product_repository_mock.go -package=mocks IProductRepository
type IProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductCoverPath(ctx context.Context, id uuid.UUID) ([]byte, error)
}

type ProductUsecase struct {
	log  *logrus.Logger
	repo IProductRepository
}

func NewProductUsecase(log *logrus.Logger, repo IProductRepository) *ProductUsecase {
	return &ProductUsecase{
		log:  log,
		repo: repo,
	}
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := u.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := u.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) GetProductCover(ctx context.Context, id uuid.UUID) ([]byte, error) {
	fileData, err := u.repo.GetProductCoverPath(ctx, id)
	if err != nil {
		return nil, err
	}

	return fileData, err
}
