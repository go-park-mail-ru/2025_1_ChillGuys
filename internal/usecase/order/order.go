package order

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IOrderUsecase interface {
	CreateOrder(context.Context, models.CreateOrderDTO) error
}

type OrderUsecase struct {
	repo order.IOrderRepository
	log  *logrus.Logger
}

func NewOrderUsecase(
	repo order.IOrderRepository,
	log *logrus.Logger,
) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		log:  log,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, in models.CreateOrderDTO) error {
	orderItems := make([]models.CreateOrderItemDTO, len(in.Items))
	for i, item := range in.Items {
		item.ID = uuid.New()
		orderItems[i] = item
	}

	order := models.Order{
		ID:                 uuid.New(),
		UserID:             in.UserID,
		Status:             models.Pending,
		TotalPrice:         20,
		TotalPriceDiscount: 10,
		AddressID:          in.AddressID,
		Items:              orderItems,
	}

	return u.repo.CreateOrder(ctx, order)
}
