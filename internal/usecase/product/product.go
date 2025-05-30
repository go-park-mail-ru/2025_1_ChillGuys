package product

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
)

//go:generate mockgen -source=product.go -destination=../../infrastructure/repository/postgres/mocks/product_repository_mock.go -package=mocks IProductRepository
type IProductRepository interface {
	GetAllProducts(ctx context.Context, offset int) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductsByCategory(
		ctx context.Context,
		id uuid.UUID,
		offset int,
		minPrice, maxPrice float64,
		minRating float32,
		sortOption models.SortOption,
	) ([]*models.Product, error)
	AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error)
	//GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Product, error)
}

type ProductUsecase struct {
	repo IProductRepository
}

func NewProductUsecase(repo IProductRepository) *ProductUsecase {
	return &ProductUsecase{
		repo: repo,
	}
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context, offset int) ([]*models.Product, error) {
	const op = "ProductUsecase.GetAllProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	products, err := u.repo.GetAllProducts(ctx, offset)
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

// GetProductsByIDs возвращает список продуктов по их UUID
func (u *ProductUsecase) GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Product, error) {
	const op = "ProductUsecase.GetProductsByIDs"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_ids", ids)

	if len(ids) == 0 {
		return []*models.Product{}, nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mu := &sync.Mutex{}
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	products := make([]*models.Product, len(ids))
	for i, id := range ids {
		wg.Add(1)
		go func(idx int, productID uuid.UUID) {
			defer wg.Done()

			if ctx.Err() != nil {
				return
			}

			product, err := u.repo.GetProductByID(ctx, productID)
			if err != nil {
				logger.WithError(err).WithField("product_id", productID).Warn("failed to get product by ID")
				trySendError(err, errCh, cancel)
				return
			}

			mu.Lock()
			products[idx] = product
			mu.Unlock()
		}(i, id)
	}

	// Горутина для закрытия канала после завершения всех операций
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Возвращаем первую ошибку (если есть)
	if err := <-errCh; err != nil {
		return nil, err
	}

	// Фильтруем nil значения (продукты, которые не удалось получить)
	filteredProducts := make([]*models.Product, 0, len(products))
	for _, p := range products {
		if p != nil {
			filteredProducts = append(filteredProducts, p)
		}
	}

	return filteredProducts, nil
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

func (u *ProductUsecase) GetProductsByCategory(
	ctx context.Context,
	id uuid.UUID,
	offset int,
	minPrice, maxPrice float64,
	minRating float32,
	sortOption models.SortOption,
) ([]*models.Product, error) {
	const op = "ProductUsecase.GetProductsByCategoryWithFilterAndSort"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("category_id", id)

	products, err := u.repo.GetProductsByCategory(
		ctx,
		id,
		offset,
		minPrice,
		maxPrice,
		minRating,
		sortOption,
	)
	if err != nil {
		logger.WithError(err).Error("get products by category with filter and sort from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	switch sortOption {
	case models.SortByPriceAsc:
		sort.Slice(products, func(i, j int) bool {
			priceI := helpers.GetFinalPrice(products[i])
			priceJ := helpers.GetFinalPrice(products[j])
			return priceI < priceJ
		})
	case models.SortByPriceDesc:
		sort.Slice(products, func(i, j int) bool {
			priceI := helpers.GetFinalPrice(products[i])
			priceJ := helpers.GetFinalPrice(products[j])
			return priceI > priceJ
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

func (u *ProductUsecase) AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error) {
	const op = "ProductUsecase.AddProduct"
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
	if product.Quantity < 0 {
		logger.Error("invalid product quantity")
		return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidProductQuantity)
	}

	// Если рейтинг не указан, ставим 0
	if product.Rating == 0 {
		product.Rating = 0
	}

	// Если количество отзывов не указано, ставим 0
	if product.ReviewsCount == 0 {
		product.ReviewsCount = 0
	}

	// Если URL превью не указан, ставим дефолтный
	if product.PreviewImageURL == "" {
		product.PreviewImageURL = "media/product-default"
	}

	// Добавляем продукт в репозиторий
	newProduct, err := u.repo.AddProduct(ctx, product, categoryID)
	if err != nil {
		logger.WithError(err).Error("add product to repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return newProduct, nil
}
