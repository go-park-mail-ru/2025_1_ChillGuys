package basket

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=basket.go -destination=../../infrastructure/repository/postgres/mocks/basket_repository_mock.go -package=mocks IBasketRepository
type IBasketRepository interface{
	Get(ctx context.Context, userID uuid.UUID) ([]*models.BasketItem, error)
	Add(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*models.BasketItem, error)
	Delete(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	UpdateQuantity(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.BasketItem, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type BasketUsecase struct {
	log  *logrus.Logger
	repo IBasketRepository
}

func NewBasketUsecase(log *logrus.Logger, repo IBasketRepository) *BasketUsecase {
	return &BasketUsecase{
		log:  log,
		repo: repo,
	}
}

func (u *BasketUsecase)Get(ctx context.Context)([]*models.BasketItem, error){
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	items, err := u.repo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}


func (u *BasketUsecase)Add(ctx context.Context, productID uuid.UUID)(*models.BasketItem, error){
	if productID == uuid.Nil {
		return nil, errs.ErrInvalidID
	}

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, errs.ErrInvalidID
	}

	item, err := u.repo.Add(ctx, userID, productID)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *BasketUsecase)Delete(ctx context.Context, productID uuid.UUID)(error){
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = u.repo.Delete(ctx, userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (u *BasketUsecase)UpdateQuantity(ctx context.Context, productID uuid.UUID, quantity int)(*models.BasketItem, error){
	if quantity <= 0 {
		return nil, errs.NewBusinessLogicError("invalid quantity")
	}

	if productID == uuid.Nil {
		return nil, errs.ErrInvalidID
	}

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	item, err := u.repo.UpdateQuantity(ctx, userID, productID, quantity)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (u *BasketUsecase)Clear(ctx context.Context,)(error){
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = u.repo.Clear(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}