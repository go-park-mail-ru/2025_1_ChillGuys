package category

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

//go:generate mockgen -source=product.go -destination=../../infrastructure/repository/postgres/mocks/product_repository_mock.go -package=mocks IProductRepository
type ICategoryRepository interface {
	GetAllCategories(ctx context.Context)([]*models.Category, error)
}

type CategoryUsecase struct {
	repo ICategoryRepository
}

func NewCategoryUsecase(repo ICategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
	}
}

func (u *CategoryUsecase) GetAllCategories(ctx context.Context)([]*models.Category, error) {
	const op = "CategoryUsecase.GetAllCategories"
    logger := logctx.GetLogger(ctx).WithField("op", op)
	
	categories, err := u.repo.GetAllCategories(ctx)
	if err != nil {
        logger.WithError(err).Error("get categories from repository")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return categories, nil
}