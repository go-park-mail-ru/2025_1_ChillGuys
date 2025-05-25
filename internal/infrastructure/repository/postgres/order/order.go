package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/google/uuid"
)

const (
	queryCreateOrder           = `INSERT INTO bazaar.order (id, user_id, status, total_price, total_price_discount, address_id) VALUES ($1, $2, $3, $4, $5, $6)`
	queryAddOrderItem          = `INSERT INTO bazaar.order_item (id, order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4, $5)`
	queryGetProductPrice       = `SELECT price, status, quantity FROM bazaar.product WHERE id = $1 LIMIT 1`
	queryGetProductDiscount    = `SELECT discounted_price, start_date, end_date FROM bazaar.discount WHERE product_id = $1`
	queryUpdateProductQuantity = `UPDATE bazaar.product SET quantity = $1 WHERE id = $2`
	queryGetOrdersByUserID     = `SELECT id, status, total_price, total_price_discount, address_id, expected_delivery_at, actual_delivery_at, created_at FROM bazaar.order WHERE user_id = $1`
	queryGetOrderProducts = `
        SELECT oi.product_id, oi.quantity, p.name 
        FROM bazaar.order_item oi
        JOIN bazaar.product p ON oi.product_id = p.id
        WHERE oi.order_id = $1`
	queryGetProductImg         = `SELECT preview_image_url FROM bazaar.product WHERE id = $1 LIMIT 1`
	queryGetOrderAddress       = `
        SELECT 
            a.region, 
            a.city, 
            a.address_string, 
            a.coordinate,
            ua.label
        FROM 
            bazaar.address a
        LEFT JOIN 
            bazaar.user_address ua ON a.id = ua.address_id
        WHERE 
            a.id = $1 
        LIMIT 1`

	queryGetOrders  = `
		SELECT id, status, total_price, total_price_discount, 
			address_id, expected_delivery_at, actual_delivery_at, created_at 
		FROM bazaar.order WHERE status = 'placed'`

	queryUpdateOrderStatus = `
		UPDATE bazaar.order
		SET 
			status = $1,
			updated_at = now()
		WHERE id = $2`

	queryGetUserIDByOrderID = `SELECT user_id FROM bazaar.order WHERE id = $1`
)

//go:generate mockgen -source=order.go -destination=../mocks/order_repository_mock.go -package=mocks IOrderRepository
type IOrderRepository interface {
	CreateOrder(context.Context, dto.CreateOrderRepoReq) error
	ProductPrice(context.Context, uuid.UUID) (*models.Product, error)
	ProductDiscounts(context.Context, uuid.UUID) ([]models.ProductDiscount, error)
	UpdateProductQuantity(context.Context, uuid.UUID, uint) error
	GetOrdersByUserID(context.Context, uuid.UUID) (*[]dto.GetOrderByUserIDResDTO, error)
	GetOrderProducts(context.Context, uuid.UUID) (*[]dto.GetOrderProductResDTO, error)
	GetProductImage(context.Context, uuid.UUID) (string, error)
	GetOrderAddress(context.Context, uuid.UUID) (*models.AddressDB, error)
	GetOrdersPlaced(ctx context.Context) (*[]dto.GetOrderByUserIDResDTO, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	GetUserIDByOrderID(ctx context.Context, orderID uuid.UUID) (uuid.UUID, error)
}

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, in dto.CreateOrderRepoReq) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	if _, err = tx.ExecContext(ctx, queryCreateOrder,
		in.Order.ID,
		in.Order.UserID,
		in.Order.Status.String(),
		in.Order.TotalPrice,
		in.Order.TotalPriceDiscount,
		in.Order.AddressID,
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("%s", err)
	}

	for productID, updatedQuantity := range in.UpdatedQuantities {
		if err = r.UpdateProductQuantity(ctx, productID, updatedQuantity); err != nil {
			tx.Rollback()
			return fmt.Errorf("%s", err)
		}
	}

	for _, item := range in.Order.Items {
		if _, err = tx.ExecContext(ctx, queryAddOrderItem,
			item.ID, in.Order.ID, item.ProductID, item.Price, item.Quantity,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("%s", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s", err)
	}

	return nil
}

func (r *OrderRepository) ProductPrice(ctx context.Context, ProductID uuid.UUID) (*models.Product, error) {
	const op = "OrderRepository.ProductPrice"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var product models.Product
	var productStatusString string

	if err := r.db.QueryRowContext(ctx, queryGetProductPrice, ProductID).Scan(
		&product.Price,
		&productStatusString,
		&product.Quantity,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("product not found")
		}
		logger.WithError(err).WithField("product_id", ProductID).Error("get product price")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	productStatus, err := models.ParseProductStatus(productStatusString)
	if err != nil {
		logger.WithError(err).WithField("status", productStatusString).Error("parse product status")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	product.Status = productStatus

	return &product, nil
}

func (r *OrderRepository) ProductDiscounts(ctx context.Context, productID uuid.UUID) ([]models.ProductDiscount, error) {
	const op = "OrderRepository.ProductDiscounts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetProductDiscount, productID)
	if err != nil {
		logger.WithError(err).WithField("product_id", productID).Error("query product discounts")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var discounts []models.ProductDiscount
	for rows.Next() {
		var discount models.ProductDiscount
		if err = rows.Scan(
			&discount.DiscountedPrice,
			&discount.DiscountStartDate,
			&discount.DiscountEndDate,
		); err != nil {
			logger.WithError(err).Error("scan discount row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		discounts = append(discounts, discount)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(discounts) == 0 {
		return nil, errs.NewNotFoundError("product discounts not found")
	}

	return discounts, nil
}

func (r *OrderRepository) UpdateProductQuantity(ctx context.Context, productID uuid.UUID, quantity uint) error {
	const op = "OrderRepository.UpdateProductQuantity"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	res, err := r.db.ExecContext(ctx, queryUpdateProductQuantity, quantity, productID)
	if err != nil {
		logger.WithError(err).WithField("product_id", productID).Error("update product quantity")
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.WithError(err).Error("get rows affected")
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return errs.NewNotFoundError("product not found for quantity update")
	}

	return nil
}

func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) (*[]dto.GetOrderByUserIDResDTO, error) {
	rows, err := r.db.QueryContext(ctx, queryGetOrdersByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []dto.GetOrderByUserIDResDTO
	for rows.Next() {
		var order dto.GetOrderByUserIDResDTO
		var status string
		if err = rows.Scan(
			&order.ID,
			&status,
			&order.TotalPrice,
			&order.TotalPriceDiscount,
			&order.AddressID,
			&order.ExpectedDeliveryAt,
			&order.ActualDeliveryAt,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orderStatus, err := models.ParseOrderStatus(status)
		if err != nil {
			return nil, err
		}

		order.Status = orderStatus
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errs.NewNotFoundError("orders not found for user")
	}

	return &orders, nil
}

func (r *OrderRepository) GetOrderProducts(ctx context.Context, orderID uuid.UUID) (*[]dto.GetOrderProductResDTO, error) {
	rows, err := r.db.QueryContext(ctx, queryGetOrderProducts, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.GetOrderProductResDTO
	for rows.Next() {
		var response dto.GetOrderProductResDTO
		if err = rows.Scan(
			&response.ProductID,
			&response.Quantity,
			&response.ProductName,
		); err != nil {
			return nil, err
		}
		result = append(result, response)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errs.NewNotFoundError("order products not found")
	}

	return &result, nil
}

func (r *OrderRepository) GetProductImage(ctx context.Context, productID uuid.UUID) (string, error) {
	var imageURL string

	if err := r.db.QueryRowContext(ctx, queryGetProductImg, productID).Scan(
		&imageURL,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.NewNotFoundError("product image not found")
		}
		return "", err
	}

	return imageURL, nil
}

func (r *OrderRepository) GetOrderAddress(ctx context.Context, addressID uuid.UUID) (*models.AddressDB, error) {
	var address models.AddressDB

	if err := r.db.QueryRowContext(ctx, queryGetOrderAddress, addressID).Scan(
		&address.Region,
		&address.City,
		&address.AddressString,
		&address.Coordinate,
		&address.Label,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("order address not found")
		}
		return nil, err
	}

	return &address, nil
}

func (w *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	const op = "AdminRepository.UpdateProductStatus"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_id", orderID)

	_, err := w.db.ExecContext(ctx, queryUpdateOrderStatus, status.String(), orderID)
	if err != nil {
		logger.WithError(err).Error("update product status")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (w *OrderRepository) GetOrdersPlaced(ctx context.Context) (*[]dto.GetOrderByUserIDResDTO, error) {
	const op = "WarehouseRepository.Get"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	var orders []dto.GetOrderByUserIDResDTO

	rows, err := w.db.QueryContext(ctx, queryGetOrders)
	if err != nil {
		logger.WithError(err).Error("query all products")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var order dto.GetOrderByUserIDResDTO
		var status string
		if err = rows.Scan(
			&order.ID,
			&status,
			&order.TotalPrice,
			&order.TotalPriceDiscount,
			&order.AddressID,
			&order.ExpectedDeliveryAt,
			&order.ActualDeliveryAt,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orderStatus, err := models.ParseOrderStatus(status)
		if err != nil {
			return nil, err
		}

		order.Status = orderStatus
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

func (r *OrderRepository) GetUserIDByOrderID(ctx context.Context, orderID uuid.UUID) (uuid.UUID, error) {
    const op = "OrderRepository.GetUserIDByOrderID"
    logger := logctx.GetLogger(ctx).WithField("op", op).WithField("order_id", orderID)

    var userID uuid.UUID

    err := r.db.QueryRowContext(ctx, queryGetUserIDByOrderID, orderID).Scan(&userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("order not found")
            return uuid.Nil, errs.NewNotFoundError("order not found")
        }
        logger.WithError(err).Error("failed to get user_id by order_id")
        return uuid.Nil, fmt.Errorf("%s: %w", op, err)
    }

    return userID, nil
}