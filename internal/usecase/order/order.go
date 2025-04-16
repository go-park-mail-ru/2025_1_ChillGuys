package order

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

//go:generate mockgen -source=order.go -destination=../mocks/order_usecase_mock.go -package=mocks IOrderUsecase
type IOrderUsecase interface {
	CreateOrder(context.Context, dto.CreateOrderDTO) error
	GetUserOrders(context.Context, uuid.UUID) (*[]dto.OrderPreviewDTO, error)
}

type OrderUsecase struct {
	repo order.IOrderRepository
}

func NewOrderUsecase(
	repo order.IOrderRepository,
) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, in dto.CreateOrderDTO) error {
	const op = "OrderUsecase.CreateOrder"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", in.UserID)
	
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
					logger.WithError(productErr).
						WithField("product_id", item.ProductID).
						Error("failed to fetch product price")
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
				if discountErr != nil && !errors.Is(discountErr, errs.ErrNotFound) {
					logger.WithError(discountErr).
						WithField("product_id", item.ProductID).
						Error("failed to fetch product discount")
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
				logger.WithFields(map[string]interface{}{
					"product_id":      item.ProductID,
					"status":          product.Status,
					"required_status": models.ProductApproved,
				}).Warn("product not approved")
				trySendError(errs.ErrProductNotApproved, errCh, cancel)
				return
			}
			if product.Quantity < item.Quantity {
				logger.WithFields(map[string]interface{}{
					"product_id":       item.ProductID,
					"requested":        item.Quantity,
					"available":        product.Quantity,
				}).Warn("not enough stock")
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
		Status:             models.Placed,
		TotalPrice:         totalPrice,
		TotalPriceDiscount: totalDiscountedPrice,
		AddressID:          in.AddressID,
		Items:              orderItems,
	}

	return u.repo.CreateOrder(ctx, dto.CreateOrderRepoReq{
		Order:             order,
		UpdatedQuantities: newQuantities,
	})
}

func (u *OrderUsecase) GetUserOrders(ctx context.Context, userID uuid.UUID) (*[]dto.OrderPreviewDTO, error) {
	const op = "OrderUsecase.GetUserOrders"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	orders, err := u.repo.GetOrdersByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("no orders found for user")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("failed to get orders by user ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	mu := &sync.Mutex{}
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	ordersPreview := make([]dto.OrderPreviewDTO, len(*orders))
	for i, orderItem := range *orders {
		wg.Add(1)
		go func() {
			defer wg.Done()

			innerWg := sync.WaitGroup{}
			var (
				address  *models.AddressDB
				products []models.OrderPreviewProductDTO
			)

			innerWg.Add(2)
			go func() {
				defer innerWg.Done()
				if ctx.Err() != nil {
					return
				}

				productIDs, productErr := u.repo.GetOrderProducts(ctx, orderItem.ID)
				if productErr != nil {
					logger.WithError(productErr).
						WithField("order_id", orderItem.ID).
						Error("failed to get order products")
					trySendError(productErr, errCh, cancel)
					return
				}

				// Получаем изображения продуктов
				productsData := make([]models.OrderPreviewProductDTO, len(*productIDs))
				imgMu := &sync.Mutex{}
				imageWg := sync.WaitGroup{}
				for i, productData := range *productIDs {
					imageWg.Add(1)

					go func() {
						defer imageWg.Done()
						if ctx.Err() != nil {
							return
						}

						productImg, imgErr := u.repo.GetProductImage(ctx, productData.ProductID)

						imgMu.Lock()
						if imgErr != nil {
							// Ошибка получения изображения, значит будем отдавать nil
							productsData[i] = models.OrderPreviewProductDTO{
								ProductImageURL: null.String{},
								ProductQuantity: productData.Quantity,
							}
							return
						}

						productsData[i] = models.OrderPreviewProductDTO{
							ProductImageURL: null.StringFrom(productImg),
							ProductQuantity: productData.Quantity,
						}
						imgMu.Unlock()
					}()
				}

				imageWg.Wait()
				products = productsData
			}()

			go func() {
				defer innerWg.Done()
				if ctx.Err() != nil {
					return
				}

				addressRes, addressErr := u.repo.GetOrderAddress(ctx, orderItem.AddressID)
				if addressErr != nil {
					logger.WithError(addressErr).
						WithField("address_id", orderItem.AddressID).
						Error("failed to get order address")
					trySendError(addressErr, errCh, cancel)
					return
				}

				address = addressRes
				address.ID = orderItem.AddressID
			}()

			innerWg.Wait()
			if ctx.Err() != nil || address == nil {
				return
			}

			mu.Lock()
			ordersPreview[i] = orderItem.ConvertToGetOrderByUserIDResDTO(address, products)
			mu.Unlock()
		}()
	}

	// Горутина для закрытия канала после завершения всех операций
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Возвращаем первую ошибку (если есть)
	if err = <-errCh; err != nil {
		logger.WithError(err).Error("failed to get order details")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
	return &ordersPreview, nil
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
