package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	queryCreateOrder           = `INSERT INTO "order" (id, user_id, status, total_price, total_price_discount, address_id) VALUES ($1, $2, $3, $4, $5, $6)`
	queryAddOrderItem          = `INSERT INTO "order_item" (id, order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4, $5)`
	queryGetProductPrice       = `SELECT price, status, quantity FROM product WHERE id = $1`
	queryGetProductDiscount    = `SELECT discounted_price, start_date, end_date FROM discount WHERE product_id = $1`
	queryUpdateProductQuantity = `UPDATE product SET quantity = $1 WHERE id = $2`
)

type IOrderRepository interface {
	CreateOrder(context.Context, models.CreateOrderRepoReq) error
	ProductPrice(context.Context, uuid.UUID) (*models.Product, error)
	ProductDiscounts(context.Context, uuid.UUID) ([]models.ProductDiscount, error)
	UpdateProductQuantity(context.Context, uuid.UUID, uint) error
}

type OrderRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewOrderRepository(db *sql.DB, log *logrus.Logger) *OrderRepository {
	return &OrderRepository{
		db:  db,
		log: log,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, in models.CreateOrderRepoReq) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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
		return err
	}

	// Обновляем количество товаров в наличии
	for productID, updatedQuantity := range in.UpdatedQuantities {
		if err = r.UpdateProductQuantity(ctx, productID, updatedQuantity); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Добавляем товары заказа
	for _, item := range in.Order.Items {
		if _, err = tx.ExecContext(ctx, queryAddOrderItem,
			item.ID, in.Order.ID, item.ProductID, item.Price, item.Quantity,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) ProductPrice(ctx context.Context, ProductID uuid.UUID) (*models.Product, error) {
	var product models.Product
	var productStatusString string

	if err := r.db.QueryRowContext(ctx, queryGetProductPrice, ProductID).Scan(
		&product.Price,
		&productStatusString,
		&product.Quantity,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	productStatus, err := models.ParseProductStatus(productStatusString)
	if err != nil {
		return nil, err
	}
	product.Status = productStatus

	return &product, nil
}

func (r *OrderRepository) ProductDiscounts(ctx context.Context, productID uuid.UUID) ([]models.ProductDiscount, error) {
	rows, err := r.db.QueryContext(ctx, queryGetProductDiscount, productID)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		discounts = append(discounts, discount)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(discounts) == 0 {
		return nil, errs.ErrNotFound
	}

	return discounts, nil
}

func (r *OrderRepository) UpdateProductQuantity(ctx context.Context, productID uuid.UUID, quantity uint) error {
	res, err := r.db.ExecContext(ctx, queryUpdateProductQuantity, quantity, productID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
