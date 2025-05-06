package search

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/guregu/null"
)

const (
	queryGetCategoryByName = `
	SELECT id, name FROM bazaar.subcategory
	WHERE LOWER(name) = LOWER($1)`
	querySearchProductsByNameWithFilterAndSort = `
SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
       p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
       d.discounted_price
FROM bazaar.product p
JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
LEFT JOIN bazaar.discount d ON p.id = d.product_id
WHERE p.status = 'approved'
  AND LOWER(p.name) LIKE LOWER($1)
  AND ($2 = '' OR ps.subcategory_id = $2::uuid)
  AND ($3 = 0 OR p.price >= $3)
  AND ($4 = 0 OR p.price <= $4)
  AND ($5 = 0::FLOAT OR p.rating >= $5::FLOAT)
ORDER BY %s
LIMIT 20 OFFSET $6`
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

func (s *SearchRepository) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	const op = "SearchRepository.GetCategoryByName"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var category models.Category

	if err := s.db.QueryRowContext(ctx, queryGetCategoryByName, name).Scan(&category.ID, &category.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.WithError(err).Error("query get category by name")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &category, nil
}

func (s *SearchRepository) GetProductsByNameWithFilterAndSort(
	ctx context.Context,
	name string,
	categoryID null.String,
	offset int,
	minPrice, maxPrice float64,
	minRating float32,
	sortOption models.SortOption,
) ([]*models.Product, error) {
	const op = "SearchRepository.GetProductsByNameWithFilterAndSort"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var orderBy string
	switch sortOption {
	case models.SortByPriceAsc:
		orderBy = "p.price ASC"
	case models.SortByPriceDesc:
		orderBy = "p.price DESC"
	case models.SortByRatingAsc:
		orderBy = "p.rating ASC"
	case models.SortByRatingDesc:
		orderBy = "p.rating DESC"
	default:
		orderBy = "p.updated_at DESC"
	}

	// Формируем запрос
	query := fmt.Sprintf(querySearchProductsByNameWithFilterAndSort, orderBy)

	// Готовим параметры запроса
	args := []interface{}{
		fmt.Sprintf("%%%s%%", name), // $1
	}

	// Добавляем параметры фильтрации по категории, цене и рейтингу
	if categoryID.Valid {
		args = append(args, categoryID.String) // $2
	} else {
		args = append(args, "") // Пустая строка, если категория не задана
	}
	args = append(args, minPrice, maxPrice, minRating, offset)

	// Выполняем запрос
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.WithError(err).Error("query search products with filter and sort")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Чтение данных
	productsList := []*models.Product{}
	for rows.Next() {
		var priceDiscount sql.NullFloat64
		product := &models.Product{}
		if err := rows.Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.PreviewImageURL,
			&product.Description,
			&product.Status,
			&product.Price,
			&product.Quantity,
			&product.UpdatedAt,
			&product.Rating,
			&product.ReviewsCount,
			&priceDiscount,
		); err != nil {
			logger.WithError(err).Error("scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		product.PriceDiscount = priceDiscount.Float64
		productsList = append(productsList, product)
	}

	// Проверка ошибок
	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}
