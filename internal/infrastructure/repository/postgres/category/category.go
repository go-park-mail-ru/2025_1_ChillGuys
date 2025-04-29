package category

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryGetAllCategories = `
			SELECT id, name FROM bazaar.category
	`

	queryGetAllSubcategories = `
		SELECT id, name 
		FROM bazaar.subcategory
		WHERE category_id = $1
	`
)

type CategoryRepository struct {
	DB  *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		DB:  db,
	}
}

func (p *CategoryRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	const op = "CategoryRepository.GetAllCategories"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	categoriesList := []*models.Category{}

	rows, err := p.DB.QueryContext(ctx, queryGetAllCategories)
	if err != nil {
		logger.WithError(err).Error("query all categories")
        return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		category := &models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
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

func (p *CategoryRepository) GetAllSubcategories(ctx context.Context, category_id uuid.UUID) ([]*models.Category, error) {
	const op = "CategoryRepository.GetAllSubcategories"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	subcategoriesList := []*models.Category{}

	rows, err := p.DB.QueryContext(ctx, queryGetAllSubcategories, category_id)
	if err != nil {
		logger.WithError(err).Error("query subcategories")
        return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		category := &models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			logger.WithError(err).Error("scan subcategory row")
            return nil, fmt.Errorf("%s: %w", op, err)
		}
		subcategoriesList = append(subcategoriesList, category)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
        return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subcategoriesList, nil
}