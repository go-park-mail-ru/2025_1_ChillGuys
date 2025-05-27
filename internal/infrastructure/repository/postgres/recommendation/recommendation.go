package recommendation

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryGetSubcategoryByProduct    = `SELECT subcategory_id FROM bazaar.product_subcategory WHERE product_id = $1`
	queryGetProductIDsBySubcategory = `SELECT product_id 
		FROM bazaar.product_subcategory 
		WHERE subcategory_id = $1
		ORDER BY RANDOM()
		LIMIT 10`
)


//go:generate mockgen -source=recommendation.go -destination=../mocks/recommendation_repository_mock.go -package=mocks IRecommendationRepository
type IRecommendationRepository interface {
	GetCategoryIDsByProductID(context.Context, uuid.UUID) ([]uuid.UUID, error)
	GetProductIDsBySubcategoryID(ctx context.Context, subcategoryID uuid.UUID, count int) ([]uuid.UUID, error)
}

type RecommendationRepository struct {
	db *sql.DB
}

func NewRecommendationRepository(db *sql.DB) *RecommendationRepository {
	return &RecommendationRepository{
		db: db,
	}
}

func (r *RecommendationRepository) GetCategoryIDsByProductID(ctx context.Context, productID uuid.UUID) ([]uuid.UUID, error) {
	const op = "RecommendationRepository.GetCategoryIDsByProductID"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetSubcategoryByProduct, productID)

	if err != nil {
		logger.WithError(err).Error("query get subcategories by product id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var categoryIDsList []uuid.UUID

	for rows.Next() {
		var subcategoryID uuid.UUID

		if err = rows.Scan(
			&subcategoryID,
		); err != nil {
			logger.WithError(err).Error("scan subcategory row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if err = rows.Err(); err != nil {
			logger.WithError(err).Error("rows iteration error")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		categoryIDsList = append(categoryIDsList, subcategoryID)
	}

	return categoryIDsList, nil
}

func (r *RecommendationRepository) GetProductIDsBySubcategoryID(ctx context.Context, subcategoryID uuid.UUID, count int) ([]uuid.UUID, error) {
	const op = "RecommendationRepository.GetProductIDsBySubcategoryID"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetProductIDsBySubcategory, subcategoryID)
	if err != nil {
		logger.WithError(err).Error("query all product ids by subcategory id")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	productIDsList := make([]uuid.UUID, 0, count)

	for rows.Next() {
		var productID uuid.UUID

		if err = rows.Scan(
			&productID,
		); err != nil {
			logger.WithError(err).Error("scan product id")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if err = rows.Err(); err != nil {
			logger.WithError(err).Error("rows iteration error")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		productIDsList = append(productIDsList, productID)
	}

	return productIDsList, nil
}
