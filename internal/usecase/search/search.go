package search

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

type ISearchRepository interface {
	GetProductsByName(context.Context, string) ([]*models.Product, error)
	GetCategoryByName(context.Context, string) (*models.Category, error)
}

type SearchUsecase struct {
	repo ISearchRepository
}

func NewSearchUsecase(repo ISearchRepository) *SearchUsecase {
	return &SearchUsecase{
		repo: repo,
	}
}

// SearchProductsByName принимает запрос с несколькими названиями продуктов и возвращает найденные продукты.
func (u *SearchUsecase) SearchProductsByName(ctx context.Context, req dto.ProductNameResponse) ([]*models.Product, error) {
	const op = "SearchUsecase.SearchProductsByName"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var allProducts []*models.Product

	// Перебираем все названия продуктов в запросе
	for _, suggestion := range req.ProductNames {
		// Вызов репозитория для поиска продукта по имени
		products, err := u.repo.GetProductsByName(ctx, suggestion.Name)
		if err != nil {
			logger.WithError(err).Error("failed to search products by name")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// Добавляем найденные продукты в общий список
		allProducts = append(allProducts, products...)
	}

	return allProducts, nil
}

// SearchCategoryByName принимает запрос с несколькими названиями категорий и возвращает найденные категории.
func (u *SearchUsecase) SearchCategoryByName(ctx context.Context, req dto.CategoryNameResponse) ([]*models.Category, error) {
	const op = "SearchUsecase.SearchCategoryByName"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var allCategories []*models.Category

	// Перебираем все названия категорий в запросе
	for _, suggestion := range req.CategoriesNames {
		// Вызов репозитория для поиска категории по имени
		category, err := u.repo.GetCategoryByName(ctx, suggestion.Name)
		if err != nil {
			logger.WithError(err).Error("failed to search category by name")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if category != nil {
			// Добавляем найденную категорию в общий список
			allCategories = append(allCategories, category)
		}
	}

	return allCategories, nil
}
