package review

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
)

type IReviewUsecase interface {
	Add(ctx context.Context, req dto.AddReviewRequest) error
	Get(ctx context.Context, req dto.GetReviewRequest) ([]*models.Review, error)
}