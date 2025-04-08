package order

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type IOrderUsecase interface {
	CreateOrder(context.Context, dto.CreateOrderDTO) error
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

func (u *OrderUsecase) CreateOrder(ctx context.Context, in dto.CreateOrderDTO) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	orderItems := make([]dto.CreateOrderItemDTO, len(in.Items))
	now := time.Now()

	var totalPrice float64 = 0
	var totalDiscountedPrice float64 = 0

	newQuantities := make(map[uuid.UUID]uint)

	mu := &sync.Mutex{}
	var totalWg sync.WaitGroup
	errCh := make(chan error, 1)

	for i, item := range in.Items {
		item.ID = uuid.New()
		orderItems[i] = item

		totalWg.Add(1)
		go func(i int, item dto.CreateOrderItemDTO) {
			defer totalWg.Done()

			var innerWg sync.WaitGroup
			innerWg.Add(2)

			var (
				product     *models.Product
				productErr  error
				discounts   []models.ProductDiscount
				discountErr error
			)

			// Получаем статус количество и цену товара
			go func() {
				defer innerWg.Done()
				if ctx.Err() != nil {
					return
				}

				product, productErr = u.repo.ProductPrice(ctx, item.ProductID)
				if productErr != nil {
					u.log.WithFields(logrus.Fields{
						"product_id": item.ProductID,
						"error":      productErr,
						"action":     "get_product_price",
					}).Error("Failed to fetch product price")
					trySendError(productErr, errCh, cancel)
					return
				}
			}()

			// Получаем скидку товара, если она есть
			go func() {
				defer innerWg.Done()
				if ctx.Err() != nil {
					return
				}

				discounts, discountErr = u.repo.ProductDiscounts(ctx, item.ProductID)
				u.log.WithFields(logrus.Fields{
					"product_id": item.ProductID,
					"error":      discountErr,
					"action":     "get_product_discount",
				}).Error("Failed to fetch product discount")
				if discountErr != nil && !errors.Is(discountErr, errs.ErrNotFound) {
					trySendError(discountErr, errCh, cancel)
					return
				}
			}()

			innerWg.Wait()

			// Проверяем отмену после запросов
			if ctx.Err() != nil {
				return
			}

			if product.Status != models.ProductApproved {
				u.log.WithFields(logrus.Fields{
					"product_id":      item.ProductID,
					"status":          product.Status,
					"required_status": models.ProductApproved,
				}).Warn("Product not approved")
				trySendError(errs.ErrProductNotApproved, errCh, cancel)
				return
			}
			if product.Quantity < item.Quantity {
				trySendError(errs.ErrNotEnoughStock, errCh, cancel)
				return
			}

			discount, _ := findLatestDiscount(discounts)

			mu.Lock()

			var priceToSave float64
			totalPrice += product.Price * float64(item.Quantity)
			newQuantities[item.ProductID] = product.Quantity - item.Quantity

			// Если есть скидка и она активна
			if discount.DiscountedPrice != 0 && discount.DiscountEndDate.After(now) {
				totalDiscountedPrice += discount.DiscountedPrice * float64(item.Quantity)
				priceToSave = discount.DiscountedPrice
			} else {
				totalDiscountedPrice += product.Price * float64(item.Quantity)
				priceToSave = product.Price
			}
			orderItems[i].Price = priceToSave
			mu.Unlock()

		}(i, item)
	}

	// Горутина для закрытия канала после завершения всех операций
	go func() {
		totalWg.Wait()
		close(errCh)
	}()

	// Возвращаем первую ошибку (если есть)
	if err := <-errCh; err != nil {
		return err
	}

	order := &dto.Order{
		ID:                 uuid.New(),
		UserID:             in.UserID,
		Status:             models.Pending,
		TotalPrice:         totalPrice,
		TotalPriceDiscount: totalDiscountedPrice,
		AddressID:          in.AddressID,
		Items:              orderItems,
	}

	u.log.Infoln(totalPrice, totalDiscountedPrice)

	return u.repo.CreateOrder(ctx, dto.CreateOrderRepoReq{
		Order:             order,
		UpdatedQuantities: newQuantities,
	})
}

// trySendError Вспомогательная функция для безопасной отправки ошибки
func trySendError(err error, errCh chan<- error, cancel context.CancelFunc) {
	select {
	case errCh <- err:
		cancel()
	default:
		// Если ошибка уже есть - игнорируем (сохраняем первую)
	}
}

// findLatestDiscount Достаём последнюю созданную скидку
func findLatestDiscount(discounts []models.ProductDiscount) (models.ProductDiscount, bool) {
	if len(discounts) == 0 {
		return models.ProductDiscount{}, false
	}

	latest := discounts[0]
	for _, discount := range discounts[1:] {
		if discount.DiscountStartDate.After(latest.DiscountStartDate) {
			latest = discount
		}
	}

	return latest, true
}
