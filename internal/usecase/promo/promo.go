package promo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

//go:generate mockgen -source=promo.go -destination=../../infrastructure/repository/postgres/mocks/promo_repository_mock.go -package=mocks IPromoRepository
type IPromoRepository interface {
	Create(ctx context.Context, promo models.PromoCode) error
	GetAll(ctx context.Context, offset int) ([]*models.PromoCode, error)
	CheckPromoCode(ctx context.Context, code string) (*models.PromoCode, error)
}

type PromoUsecase struct {
	repo IPromoRepository
}

func NewPromoUsecase(repo IPromoRepository) *PromoUsecase{
	return &PromoUsecase{repo: repo}
}

func (uc *PromoUsecase) CreatePromo(ctx context.Context, req dto.CreatePromoRequest) (dto.PromoResponse, error) {
	const op = "PromoUsecase.CreatePromo"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	promoDB := models.PromoCode{
		ID:        uuid.New(),
		Code:      req.Code,
		Percent:   req.Percent,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := uc.repo.Create(ctx, promoDB); err != nil {
		logger.WithError(err).Error("failed to create promo")
		return dto.PromoResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.PromoResponse{
		ID:        promoDB.ID,
		Code:      promoDB.Code,
		Percent:   promoDB.Percent,
		StartDate: promoDB.StartDate,
		EndDate:   promoDB.EndDate,
	}, nil
}

func (uc *PromoUsecase) GetAllPromos(ctx context.Context, offset int) (dto.PromosResponse, error) {
	const op = "PromoUsecase.GetAllPromos"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	promos, err := uc.repo.GetAll(ctx, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get promos")
		return dto.PromosResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	if promos == nil {
		return dto.PromosResponse{}, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
	}

	return dto.ConvertToPromosResponse(promos), nil
}

func (uc *PromoUsecase) CheckPromoCode(ctx context.Context, code string) (dto.PromoValidityResponse, error) {
    const op = "PromoUsecase.GetAllPromos"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	promo, err := uc.repo.CheckPromoCode(ctx, code)
    if err != nil {
		logger.WithError(err).Error("failed to get promo")
        return dto.PromoValidityResponse{IsValid: false}, nil
    }
    
    now := time.Now()
    isValid := now.After(promo.StartDate) && now.Before(promo.EndDate)
    
    response := dto.PromoValidityResponse{
        IsValid: isValid,
    }
    
    if isValid {
        response.Percent = promo.Percent
    }
    
    return response, nil
}