package basket

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=basket.go -destination=../../infrastructure/repository/postgres/mocks/basket_repository_mock.go -package=mocks IBasketRepository
type IBasketRepository interface{
	GetOrCreateBasket(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	GetProductsInBasket(ctx context.Context, userID uuid.UUID) (*dto.BasketResponse, error)
	AddProductInBasket(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*models.BasketItem, error)
	DeleteProductFromBasket(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
	UpdateProductQuantity(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.BasketItem, error)
	ClearBasket(ctx context.Context, userID uuid.UUID) error
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

func (u *BasketUsecase)Get(ctx context.Context)(*dto.BasketResponse, error){
	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	items, err := u.repo.GetProductsInBasket(ctx, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}


func (u *BasketUsecase)AddProduct(ctx context.Context, productID uuid.UUID)(*models.BasketItem, error){
	if productID == uuid.Nil {
		return nil, errs.ErrInvalidProductID
	}

	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	item, err := u.repo.AddProductInBasket(ctx, userID, productID)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *BasketUsecase)DeleteProduct(ctx context.Context, productID uuid.UUID)(error){
	if productID == uuid.Nil {
		return errs.ErrInvalidProductID
	}

	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = u.repo.DeleteProductFromBasket(ctx, userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (u *BasketUsecase)UpdateProductQuantity(ctx context.Context, productID uuid.UUID, quantity int)(*models.BasketItem, error){
	if quantity <= 0 {
		return nil, errs.ErrInvalidQuantity
	}

	if productID == uuid.Nil {
		return nil, errs.ErrInvalidProductID
	}

	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	item, err := u.repo.UpdateProductQuantity(ctx, userID, productID, quantity)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (u *BasketUsecase)Clear(ctx context.Context,)(error){
	userID, err := u.getUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = u.repo.ClearBasket(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *BasketUsecase) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
    userIDStr, isExist := ctx.Value(cookie.UserIDKey).(string)
    if !isExist {
        return uuid.Nil, errs.ErrUserNotFound
    }
    
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.Nil, errs.ErrInvalidUserID
    }
    
    return userID, nil
}