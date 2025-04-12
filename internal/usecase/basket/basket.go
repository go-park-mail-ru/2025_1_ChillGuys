package basket

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/google/uuid"
)

//go:generate mockgen -source=basket.go -destination=../../infrastructure/repository/postgres/mocks/basket_repository_mock.go -package=mocks IBasketRepository
type IBasketRepository interface{
	Get(ctx context.Context, userID uuid.UUID) ([]*models.BasketItem, error)
	Add(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*models.BasketItem, error)
	Delete(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	UpdateQuantity(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.BasketItem, int, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type BasketUsecase struct {
	repo IBasketRepository
}

func NewBasketUsecase(repo IBasketRepository) *BasketUsecase {
	return &BasketUsecase{
		repo: repo,
	}
}

func (u *BasketUsecase)Get(ctx context.Context)([]*models.BasketItem, error){
	const op = "BasketUsecase.Get"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	logger = logger.WithField("user_id", userID)
    items, err := u.repo.Get(ctx, userID)
    if err != nil {
        logger.WithError(err).Error("get basket items from repo")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return items, nil
}


func (u *BasketUsecase)Add(ctx context.Context, productID uuid.UUID)(*models.BasketItem, error){
	const op = "BasketUsecase.Add"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    if productID == uuid.Nil {
        logger.Error("invalid product ID")
        return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
    }

	userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
    }

	logger.WithField("user_id", userID).WithField("product_id", productID)

	item, err := u.repo.Add(ctx, userID, productID)
    if err != nil {
        logger.WithError(err).Error("add product to basket")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return item, nil
}

func (u *BasketUsecase)Delete(ctx context.Context, productID uuid.UUID)(error){
	const op = "BasketUsecase.Delete"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return fmt.Errorf("%s: %w", op, err)
    }

	logger.WithField("user_id", userID).WithField("product_id", productID)

	err = u.repo.Delete(ctx, userID, productID)
	if err != nil {
        logger.WithError(err).Error("delete product from basket")
        return fmt.Errorf("%s: %w", op, err)
    }

	return nil
}

func (u *BasketUsecase)UpdateQuantity(ctx context.Context, productID uuid.UUID, quantity int)(*models.BasketItem, int, error){
	const op = "BasketUsecase.UpdateQuantity"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    if quantity <= 0 {
        logger.WithField("quantity", quantity).Error("invalid quantity")
        return nil, -1, fmt.Errorf("%s: %w", op, errs.NewBusinessLogicError("invalid quantity"))
    }

	if productID == uuid.Nil {
        logger.Error("invalid product ID")
        return nil, -1, fmt.Errorf("%s: %w", op, errs.ErrInvalidID)
    }

	userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return nil, -1, fmt.Errorf("%s: %w", op, err)
    }

	logger.WithField("user_id", userID).WithField("product_id", productID)

	item, rem ,err := u.repo.UpdateQuantity(ctx, userID, productID, quantity)
    if err != nil {
        logger.WithError(err).Error("update product quantity")
        return nil, -1, fmt.Errorf("%s: %w", op, err)
    }

	return item, rem, nil
}

func (u *BasketUsecase)Clear(ctx context.Context,)(error){
	const op = "BasketUsecase.Clear"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    userID, err := helpers.GetUserIDFromContext(ctx)
    if err != nil {
        logger.WithError(err).Error("get user ID from context")
        return fmt.Errorf("%s: %w", op, err)
    }

	logger = logger.WithField("user_id", userID)
    err = u.repo.Clear(ctx, userID)
    if err != nil {
        logger.WithError(err).Error("clear basket")
        return fmt.Errorf("%s: %w", op, err)
    }

	return nil
}