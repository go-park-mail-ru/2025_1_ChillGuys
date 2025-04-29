package category

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

//go:generate mockgen -source=category.go -destination=../../infrastructure/repository/postgres/mocks/category_repository_mock.go -package=mocks ICategoryRepository
type ICategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetAllSubcategories(ctx context.Context, category_id uuid.UUID) ([]*models.Category, error)
}

type CategoryUsecase struct {
	repo ICategoryRepository
}

func NewCategoryUsecase(repo ICategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		repo: repo,
	}
}

func (u *CategoryUsecase) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	const op = "CategoryUsecase.GetAllCategories"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	categories, err := u.repo.GetAllCategories(ctx)
	if err != nil {
		logger.WithError(err).Error("get categories from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return categories, nil
}

func (u *CategoryUsecase) GetAllSubategories(ctx context.Context, category_id uuid.UUID) ([]*models.Category, error) {
	const op = "CategoryUsecase.GetAllCategories"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	categories, err := u.repo.GetAllSubcategories(ctx, category_id)
	if err != nil {
		logger.WithError(err).Error("get subcategories from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return categories, nil
}
