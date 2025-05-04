package search

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/guregu/null"
	"sort"
	"sync"
)

type ISearchRepository interface {
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	GetProductsByName(ctx context.Context, name string, categoryID null.String, offset int) ([]*models.Product, error)
	GetProductsByNameWithFilterAndSort(
		ctx context.Context,
		name string,
		categoryID null.String,
		offset int,
		minPrice, maxPrice float64,
		minRating float32,
		sortOption models.SortOption,
	) ([]*models.Product, error)
}

type SearchUsecase struct {
	repo ISearchRepository
}

func NewSearchUsecase(repo ISearchRepository) *SearchUsecase {
	return &SearchUsecase{
		repo: repo,
	}
}

func (u *SearchUsecase) SearchProductsByNameWithFilterAndSort(
	ctx context.Context,
	categoryID null.String,
	subString string,
	offset int,
	minPrice, maxPrice float64,
	minRating float32,
	sortOption models.SortOption,
) ([]*models.Product, error) {
	const op = "SearchUsecase.SearchProductsByNameWithFilterAndSort"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("sub_string", subString)

	products, err := u.repo.GetProductsByNameWithFilterAndSort(ctx, subString, categoryID, offset, minPrice, maxPrice, minRating, sortOption)
	if err != nil {
		logger.WithError(err).Warn("failed to search products with filter and sort")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Локальная сортировка, если нужно (если БД не справляется)
	switch sortOption {
	case models.SortByPriceAsc:
		sort.Slice(products, func(i, j int) bool {
			return helpers.GetFinalPrice(products[i]) < helpers.GetFinalPrice(products[j])
		})
	case models.SortByPriceDesc:
		sort.Slice(products, func(i, j int) bool {
			return helpers.GetFinalPrice(products[i]) > helpers.GetFinalPrice(products[j])
		})
	case models.SortByRatingAsc:
		sort.Slice(products, func(i, j int) bool {
			return products[i].Rating < products[j].Rating
		})
	case models.SortByRatingDesc:
		sort.Slice(products, func(i, j int) bool {
			return products[i].Rating > products[j].Rating
		})
	}

	return products, nil
}

// trySendError Вспомогательная функция для безопасной отправки ошибки
func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}

func (u *SearchUsecase) SearchProductsBySubString(
	ctx context.Context,
	subString string,
	categoryID null.String,
	offset int,
) ([]*models.Product, error) {
	const op = "SearchUsecase.SearchProductsBySubString"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("sub_string", subString)

	products, err := u.repo.GetProductsByName(ctx, subString, categoryID, offset)
	if err != nil {
		logger.WithError(err).Warn("failed to search products by substring")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (u *SearchUsecase) SearchCategoryByName(ctx context.Context, req dto.CategoryNameResponse) ([]*models.Category, error) {
	const op = "SearchUsecase.SearchCategoryByName"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("names_count", len(req.CategoriesNames))

	if len(req.CategoriesNames) == 0 {
		return []*models.Category{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	errCh := make(chan error, 1)
	categories := make([]*models.Category, len(req.CategoriesNames))

	for i, suggestion := range req.CategoriesNames {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			category, err := u.repo.GetCategoryByName(ctx, name)
			if err != nil {
				logger.WithError(err).WithField("category_name", name).Warn("failed to search category by name")
				trySendError(err, errCh, cancel)
				return
			}

			if category != nil {
				mu.Lock()
				categories[idx] = category
				mu.Unlock()
			}
		}(i, suggestion.Name)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	filtered := make([]*models.Category, 0, len(categories))
	for _, c := range categories {
		if c != nil {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}
