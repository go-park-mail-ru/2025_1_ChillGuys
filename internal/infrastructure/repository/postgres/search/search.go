package search

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

const (
	querySearchProductsByName = `
    SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
        p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
        d.discounted_price
    FROM bazaar.product p 
    LEFT JOIN bazaar.discount d ON p.id = d.product_id
    WHERE p.status = 'approved' AND LOWER(p.name) LIKE LOWER($1)`
	queryGetCategoryByName = `
	SELECT id, name FROM bazaar.subcategory
	WHERE LOWER(name) = LOWER($1)`
	querySearchProductsByNameWithFilterAndSort = `
		SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
			p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
			d.discounted_price
		FROM 
			bazaar.product p 
		LEFT JOIN 
			bazaar.discount d ON p.id = d.product_id
		WHERE 
			p.status = 'approved' 
			AND LOWER(p.name) LIKE LOWER($1)
			AND ($3 = 0 OR p.price > $3)
			AND ($4 = 0 OR p.price < $4)
			AND ($5 = 0::FLOAT OR p.rating > $5::FLOAT)
		LIMIT 10 OFFSET $2`
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		db: db,
	}
}

func (s *SearchRepository) GetProductsByName(ctx context.Context, name string) ([]*models.Product, error) {
	const op = "SearchRepository.GetProductsByName"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.Product{}

	// Подготовка строки для поиска с % в начале и конце
	pattern := fmt.Sprintf("%%%s%%", name)

	// Выполнение запроса
	rows, err := s.db.QueryContext(ctx, querySearchProductsByName, pattern)
	if err != nil {
		logger.WithError(err).Error("query search products by name")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	// Чтение строк
	for rows.Next() {
		var priceDiscount sql.NullFloat64
		product := &models.Product{}
		if err = rows.Scan(
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

	// Проверка ошибок при переборе строк
	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
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
	offset int,
	minPrice, maxPrice float64,
	minRating float32,
) ([]*models.Product, error) {
	const op = "SearchRepository.GetProductsByNameWithFilterAndSort"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.Product{}
	pattern := fmt.Sprintf("%%%s%%", name)

	rows, err := s.db.QueryContext(
		ctx,
		querySearchProductsByNameWithFilterAndSort,
		pattern,
		offset,
		minPrice,
		maxPrice,
		minRating,
	)

	if err != nil {
		logger.WithError(err).Error("query search products by name with filter and sort")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var priceDiscount sql.NullFloat64
		product := &models.Product{}
		if err = rows.Scan(
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

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}
