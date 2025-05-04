package suggestions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/lib/pq"
)

const (
	queryGetAllCategoriesName      = `SELECT name FROM bazaar.subcategory`
	queryGetAllProductsName        = `SELECT name FROM bazaar.product WHERE status = 'approved'`
	queryGetProductsNameByCategory = `
		SELECT DISTINCT p.name
		FROM bazaar.product p
		JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
		JOIN bazaar.subcategory s ON s.id = ps.subcategory_id
		WHERE s.id = $1 AND p.status = 'approved'`
	queryGetAllProductsNamePaginated = `
        SELECT name 
		FROM bazaar.product 
		WHERE status = 'approved'
		ORDER BY name, id  
		LIMIT 20 OFFSET $1`
	queryGetProductsNameByCategoryPaginated = `
        SELECT DISTINCT p.name
        FROM bazaar.product p
        JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
        JOIN bazaar.subcategory s ON s.id = ps.subcategory_id
        WHERE s.id = $1 AND p.status = 'approved'
        ORDER BY p.name
        LIMIT 20 OFFSET $2`
	queryCountAllProducts = `
        SELECT COUNT(*) 
        FROM bazaar.product 
        WHERE status = 'approved'`
	queryCountProductsByCategory = `
        SELECT COUNT(DISTINCT p.id)
        FROM bazaar.product p
        JOIN bazaar.product_subcategory ps ON p.id = ps.product_id
        WHERE ps.subcategory_id = $1 AND p.status = 'approved'`
)

type SuggestionsRepository struct {
	db *sql.DB
}

func NewSuggestionsRepository(db *sql.DB) *SuggestionsRepository {
	return &SuggestionsRepository{
		db: db,
	}
}

func (p *SuggestionsRepository) GetAllCategoriesName(ctx context.Context) ([]*models.CategorySuggestion, error) {
	const op = "CategoryRepository.GetAllCategories"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	categoriesList := []*models.CategorySuggestion{}

	rows, err := p.db.QueryContext(ctx, queryGetAllCategoriesName)
	if err != nil {
		logger.WithError(err).Error("query all categories")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		category := &models.CategorySuggestion{}
		if err = rows.Scan(
			&category.Name,
		); err != nil {
			logger.WithError(err).Error("scan category row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		categoriesList = append(categoriesList, category)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return categoriesList, nil
}

func (p *SuggestionsRepository) GetAllProductsName(ctx context.Context) ([]*models.ProductSuggestion, error) {
	const op = "ProductRepository.GetAllProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.ProductSuggestion{}

	rows, err := p.db.QueryContext(ctx, queryGetAllProductsName)
	if err != nil {
		logger.WithError(err).Error("query all products")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		product := &models.ProductSuggestion{}
		if err = rows.Scan(
			&product.Name,
		); err != nil {
			logger.WithError(err).Error("scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}

func (p *SuggestionsRepository) GetProductsNameByCategory(ctx context.Context, categoryID string) ([]*models.ProductSuggestion, error) {
	const op = "ProductRepository.GetProductsNameByCategory"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.ProductSuggestion{}

	rows, err := p.db.QueryContext(ctx, queryGetProductsNameByCategory, categoryID)
	if err != nil {
		logger.WithError(err).Error("query products by category")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		product := &models.ProductSuggestion{}
		if err = rows.Scan(&product.Name); err != nil {
			logger.WithError(err).Error("scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}

func (p *SuggestionsRepository) GetAllProductsNameOffset(ctx context.Context, pageNum int) ([]*models.ProductSuggestion, error) {
	const op = "SuggestionsRepository.GetAllProductsNameOffset"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	const limit = 20

	if pageNum < 0 {
		pageNum = 0
		logger.Warn("Negative page number provided, resetting to 0")
	}

	// Получаем общее количество продуктов (уникальных имен)
	var total int
	err := p.db.QueryRowContext(ctx, `
        SELECT COUNT(DISTINCT name) 
        FROM bazaar.product 
        WHERE status = 'approved'
    `).Scan(&total)
	if err != nil {
		logger.WithError(err).Error("failed to count distinct product names")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Вычисляем максимальный номер страницы (округление вверх)
	maxPage := (total + limit - 1) / limit
	if pageNum >= maxPage {
		return []*models.ProductSuggestion{}, nil
	}

	// Вычисляем смещение
	offset := pageNum * limit

	// 1. Сначала получаем список уникальных имен с пагинацией
	uniqueNames := make([]string, 0, limit)
	rows, err := p.db.QueryContext(ctx, `
        SELECT DISTINCT name 
        FROM bazaar.product 
        WHERE status = 'approved'
        ORDER BY name
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query distinct product names")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			logger.WithError(err).Error("failed to scan product name")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		uniqueNames = append(uniqueNames, name)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(uniqueNames) == 0 {
		return []*models.ProductSuggestion{}, nil
	}

	// 2. Теперь получаем ВСЕ продукты с этими именами (включая дубликаты)
	productsList := make([]*models.ProductSuggestion, 0)
	rows, err = p.db.QueryContext(ctx, `
        SELECT name 
        FROM bazaar.product 
        WHERE status = 'approved' AND name = ANY($1)
        ORDER BY name, id
    `, pq.Array(uniqueNames)) // pq.Array - для передачи массива в PostgreSQL
	if err != nil {
		logger.WithError(err).Error("failed to query products by names")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		product := &models.ProductSuggestion{}
		if err = rows.Scan(&product.Name); err != nil {
			logger.WithError(err).Error("failed to scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}

func (p *SuggestionsRepository) GetProductsNameByCategoryOffset(ctx context.Context, categoryID string, pageNum int) ([]*models.ProductSuggestion, error) {
	const op = "SuggestionsRepository.GetProductsNameByCategoryOffset"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	const limit = 20

	if pageNum < 0 {
		pageNum = 0
		logger.Warn("Negative page number provided, resetting to 0")
	}

	// Получаем общее количество продуктов в категории
	var total int
	err := p.db.QueryRowContext(ctx, queryCountProductsByCategory, categoryID).Scan(&total)
	if err != nil {
		logger.WithError(err).Error("failed to count products in category")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Вычисляем максимальный номер страницы (округление вверх)
	maxPage := (total + limit - 1) / limit
	if pageNum >= maxPage {
		return []*models.ProductSuggestion{}, nil
	}

	// Вычисляем абсолютное смещение
	offset := pageNum * limit

	productsList := make([]*models.ProductSuggestion, 0, limit)
	rows, err := p.db.QueryContext(ctx, queryGetProductsNameByCategoryPaginated, categoryID, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query products by category with pagination")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		product := &models.ProductSuggestion{}
		if err = rows.Scan(&product.Name); err != nil {
			logger.WithError(err).Error("failed to scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}
