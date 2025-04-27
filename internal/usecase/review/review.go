package review

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

type IReviewRepository interface{
	AddReview(ctx context.Context, review models.ReviewDB) error
	GetReview(ctx context.Context, productID uuid.UUID, offset int) ([]*models.Review, error)
}

type ReviewUsecase struct {
	repo IReviewRepository
}

func NewReviewUsecase(repo IReviewRepository) *ReviewUsecase {
	return &ReviewUsecase{
		repo : repo,
	}
}

func (u *ReviewUsecase) Add(ctx context.Context, req dto.AddReviewRequest) error {
	const op = "ReviewUsecase.Add"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	userIDStr, isExist := ctx.Value(domains.UserIDKey{}).(string)
	if !isExist || userIDStr == "" {
		logger.Warn("user ID not found in context")
		return fmt.Errorf("%s: %w", op, errs.ErrNotFound)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.WithError(err).Error("invalid user ID format")
		return fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
	}
	
	review := models.ReviewDB{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := u.repo.AddReview(ctx, review); err != nil {
		return err
	}

	return nil
}

func (u *ReviewUsecase) Get(ctx context.Context, req dto.GetReviewRequest) ([]*models.Review, error) {
	reviews, err := u.repo.GetReview(ctx, req.ProductID, req.Offset)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}