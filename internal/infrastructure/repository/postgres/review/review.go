package review

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	queryAddReview = `
		INSERT INTO bazaar.review (id, user_id, product_id, rating, comment)
			VALUES ($1, $2, $3, $4, $5)
	`

	queryUpdateCount = `
		UPDATE bazaar.product
		SET reviews_count = reviews_count + 1 
		WHERE id = $1
	`

	queryGetReview = `
		SELECT 
			r.id, u.name, u.surname, u.image_url, r.rating, r.comment
		FROM bazaar.review r
		JOIN bazaar.user u ON r.user_id = u.id
		WHERE r.product_id = $1
		ORDER BY r.created_at DESC
        LIMIT 7 OFFSET $2
	`
)

type ReviewRepository struct {
	DB *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{
		DB:  db,
	}
}

func (r *ReviewRepository) AddReview(ctx context.Context, review models.ReviewDB) error {
	const op = "AddReview.AddReview"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("begin transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, queryUpdateCount, review.ProductID)
	if err != nil {
		logger.WithError(err).Error("increment count")
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.ExecContext(ctx, queryAddReview,
		review.ID,
		review.UserID,
		review.ProductID,
		review.Rating,
		review.Comment,
	)
	if err != nil {
		// Проверяем, является ли ошибка нарушением уникальности
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // Код ошибки "unique_violation"
				logger.WithError(err).Error("user has already reviewed this product")
				tx.Rollback()
				return errs.ErrAlreadyExists
			}
		}

		logger.WithError(err).Error("add review")
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		logger.WithError(err).Error("commit transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *ReviewRepository) GetReview(ctx context.Context, productID uuid.UUID, offset int) ([]*models.Review, error) {
	const op = "ReviewRepository.GetReview"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("productID", productID)

	reviewList := []*models.Review{}

	rows, err := r.DB.QueryContext(ctx, queryGetReview, productID, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("no reviews on this product")
            return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("query reviews")
        return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		review := &models.Review{}
		err = rows.Scan(
			&review.ID,
			&review.Name,
			&review.Surname,
			&review.ImageURL,
			&review.Rating,
			&review.Comment,
		)
		if err != nil {
			logger.WithError(err).Error("scan review")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reviewList = append(reviewList, review)
	}

	if err = rows.Err(); err != nil {
        logger.WithError(err).Error("rows iteration error")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return reviewList, nil
}