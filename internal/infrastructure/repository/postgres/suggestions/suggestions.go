package suggestions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

const (
	queryGetAllCategoriesName = `SELECT name FROM bazaar.subcategory`
	queryGetAllProductsName   = `SELECT name FROM bazaar.product WHERE status = 'approved'`
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
