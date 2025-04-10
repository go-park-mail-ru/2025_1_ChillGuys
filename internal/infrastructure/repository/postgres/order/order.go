package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryCreateOrder           = `INSERT INTO bazaar."order" (id, user_id, status, total_price, total_price_discount, address_id) VALUES ($1, $2, $3, $4, $5, $6)`
	queryAddOrderItem          = `INSERT INTO bazaar."order_item" (id, order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4, $5)`
	queryGetProductPrice       = `SELECT price, status, quantity FROM bazaar.product WHERE id = $1`
	queryGetProductDiscount    = `SELECT discounted_price, start_date, end_date FROM bazaar.discount WHERE product_id = $1`
	queryUpdateProductQuantity = `UPDATE bazaar.product SET quantity = $1 WHERE id = $2`
)

type IOrderRepository interface {
	CreateOrder(context.Context, models.CreateOrderRepoReq) error
	ProductPrice(context.Context, uuid.UUID) (*models.Product, error)
	ProductDiscounts(context.Context, uuid.UUID) ([]models.ProductDiscount, error)
	UpdateProductQuantity(context.Context, uuid.UUID, uint) error
}

type OrderRepository struct {
	db  *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db:  db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, in models.CreateOrderRepoReq) error {
	const op = "OrderRepository.CreateOrder"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("failed to begin transaction")
		return fmt.Errorf("%s: %w", op, err)
	}

	// Создаём заказ
	if _, err = tx.ExecContext(ctx, queryCreateOrder,
		in.Order.ID,
		in.Order.UserID,
		in.Order.Status.String(),
		in.Order.TotalPrice,
		in.Order.TotalPriceDiscount,
		in.Order.AddressID,
	); err != nil {
		tx.Rollback()
		logger.WithError(err).Error("create order")
		return fmt.Errorf("%s: %w", op, err)
	}

	// Обновляем количество товаров в наличии
	for productID, updatedQuantity := range in.UpdatedQuantities {
		if err = r.UpdateProductQuantity(ctx, productID, updatedQuantity); err != nil {
			tx.Rollback()
			logger.WithError(err).WithField("product_id", productID).Error("update product quantity")
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	// Добавляем товары заказа
	for _, item := range in.Order.Items {
		if _, err = tx.ExecContext(ctx, queryAddOrderItem,
			item.ID, in.Order.ID, item.ProductID, item.Price, item.Quantity,
		); err != nil {
			tx.Rollback()
			logger.WithError(err).WithField("product_id", item.ProductID).Error("add order item")
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("commit transaction")
		return fmt.Errorf("%s: %w", op, err)
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
			logger.WithField("product_id", ProductID).Warn("product not found")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
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
		logger.WithField("product_id", productID).Warn("no discounts found")
		return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
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
		logger.WithField("product_id", productID).Warn("product not found")
		return fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
	}

	return nil
}