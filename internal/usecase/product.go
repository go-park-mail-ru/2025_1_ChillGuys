package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
)

type ProductUsecase struct {
	log   *logrus.Logger
	repo  repository.IProductRepository
}

func NewProductUsecase(log *logrus.Logger, repo repository.IProductRepository) *ProductUsecase{
	return &ProductUsecase{
		log: log,
		repo: repo,
	}
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]*models.Product, error){
	products, err := u.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error){
	product, err := u.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductUsecase) GetProductCover(ctx context.Context, id uuid.UUID) ([]byte, error){
	fileData, err := u.repo.GetProductCoverPath(ctx, id)
	if err != nil {
		return nil, err
	}

	return fileData, err
}