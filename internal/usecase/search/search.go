package search

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"sync"
)

type ISearchRepository interface {
	GetProductsByName(context.Context, string, int) ([]*models.Product, error)
	GetCategoryByName(context.Context, string) (*models.Category, error)
	GetProductsByNameWithFilterAndSort(
		ctx context.Context, 
		name string, 
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

// SearchProductsByName принимает запрос с несколькими названиями продуктов и возвращает найденные продукты.
func (u *SearchUsecase) SearchProductsByName(ctx context.Context, req dto.ProductNameResponse, offset int) ([]*models.Product, error) {
	const op = "SearchUsecase.SearchProductsByName"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("names_count", len(req.ProductNames))

	if len(req.ProductNames) == 0 {
		return []*models.Product{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	errCh := make(chan error, 1)
	allProducts := make([][]*models.Product, len(req.ProductNames))

	for i, suggestion := range req.ProductNames {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			products, err := u.repo.GetProductsByName(ctx, name, offset)
			if err != nil {
				logger.WithError(err).WithField("product_name", name).Warn("failed to search products by name")
				trySendError(err, errCh, cancel)
				return
			}

			mu.Lock()
			allProducts[idx] = products
			mu.Unlock()
		}(i, suggestion.Name)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Объединяем все слайсы продуктов в один
	var merged []*models.Product
	for _, products := range allProducts {
		merged = append(merged, products...)
	}

	return merged, nil
}

// SearchCategoryByName принимает запрос с несколькими названиями категорий и возвращает найденные категории.
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

// trySendError Вспомогательная функция для безопасной отправки ошибки
func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}

func (u *SearchUsecase) SearchProductsByNameWithFilterAndSort(
    ctx context.Context, 
    req dto.ProductNameResponse, 
    offset int,
    minPrice, maxPrice float64,
    minRating float32,
    sortOption models.SortOption,
) ([]*models.Product, error) {
    const op = "SearchUsecase.SearchProductsByNameWithFilterAndSort"
    logger := logctx.GetLogger(ctx).WithField("op", op).WithField("names_count", len(req.ProductNames))

    if len(req.ProductNames) == 0 {
        return []*models.Product{}, nil
    }

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    var wg sync.WaitGroup
    mu := &sync.Mutex{}
    errCh := make(chan error, 1)
    allProducts := make([][]*models.Product, len(req.ProductNames))

    for i, suggestion := range req.ProductNames {
        wg.Add(1)
        go func(idx int, name string) {
            defer wg.Done()

            if ctx.Err() != nil {
                return
            }

            products, err := u.repo.GetProductsByNameWithFilterAndSort(
                ctx, 
                name, 
                offset,
                minPrice,
                maxPrice,
                minRating,
                sortOption,
            )
            if err != nil {
                logger.WithError(err).WithField("product_name", name).Warn("failed to search products by name")
                trySendError(err, errCh, cancel)
                return
            }

            mu.Lock()
            allProducts[idx] = products
            mu.Unlock()
        }(i, suggestion.Name)
    }

    go func() {
        wg.Wait()
        close(errCh)
    }()

    if err := <-errCh; err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    var merged []*models.Product
    for _, products := range allProducts {
        merged = append(merged, products...)
    }

    return merged, nil
}