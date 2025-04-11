package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	queryCreateOrder           = `INSERT INTO bazaar."order" (id, user_id, status, total_price, total_price_discount, address_id) VALUES ($1, $2, $3, $4, $5, $6)`
	queryAddOrderItem          = `INSERT INTO bazaar."order_item" (id, order_id, product_id, price, quantity) VALUES ($1, $2, $3, $4, $5)`
	queryGetProductPrice       = `SELECT price, status, quantity FROM bazaar.product WHERE id = $1 LIMIT 1`
	queryGetProductDiscount    = `SELECT discounted_price, start_date, end_date FROM bazaar.discount WHERE product_id = $1`
	queryUpdateProductQuantity = `UPDATE bazaar.product SET quantity = $1 WHERE id = $2`
	queryGetOrdersByUserID     = `SELECT id, status, total_price, total_price_discount, address_id, expected_delivery_at, actual_delivery_at, created_at FROM bazaar."order" WHERE user_id = $1`
	queryGetOrderProducts      = `SELECT product_id, quantity FROM bazaar.order_item WHERE order_id = $1`
	queryGetProductImg         = `SELECT preview_image_url FROM bazaar.product WHERE id = $1 LIMIT 1`
	queryGetOrderAddress       = `SELECT city, street, house, apartment, zip_code FROM bazaar.address WHERE id = $1 LIMIT 1`
)

type IOrderRepository interface {
	CreateOrder(context.Context, dto.CreateOrderRepoReq) error
	ProductPrice(context.Context, uuid.UUID) (*models.Product, error)
	ProductDiscounts(context.Context, uuid.UUID) ([]models.ProductDiscount, error)
	UpdateProductQuantity(context.Context, uuid.UUID, uint) error
	GetOrdersByUserID(context.Context, uuid.UUID) (*[]dto.GetOrderByUserIDResDTO, error)
	GetOrderProducts(context.Context, uuid.UUID) (*[]dto.GetOrderProductResDTO, error)
	GetProductImage(context.Context, uuid.UUID) (string, error)
	GetOrderAddress(context.Context, uuid.UUID) (*models.AddressDB, error)
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

func (r *OrderRepository) CreateOrder(ctx context.Context, in dto.CreateOrderRepoReq) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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
		return err
	}

	for productID, updatedQuantity := range in.UpdatedQuantities {
		if err = r.UpdateProductQuantity(ctx, productID, updatedQuantity); err != nil {
			tx.Rollback()
			return err
		}
	}

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
		return nil, errs.ErrNotFound
	}

	return &orders, nil
}

func (r *OrderRepository) GetOrderProducts(ctx context.Context, productID uuid.UUID) (*[]dto.GetOrderProductResDTO, error) {
	rows, err := r.db.QueryContext(ctx, queryGetOrderProducts, productID)
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
		); err != nil {
			return nil, err
		}
		result = append(result, response)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errs.ErrNotFound
	}

	return &result, nil
}

func (r *OrderRepository) GetProductImage(ctx context.Context, productID uuid.UUID) (string, error) {
	var imageURL string

	if err := r.db.QueryRowContext(ctx, queryGetProductImg, productID).Scan(
		&imageURL,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrNotFound
		}
		return "", err
	}

	return imageURL, nil
}

func (r *OrderRepository) GetOrderAddress(ctx context.Context, addressID uuid.UUID) (*models.AddressDB, error) {
	var address models.AddressDB

	if err := r.db.QueryRowContext(ctx, queryGetOrderAddress, addressID).Scan(
		&address.City,
		&address.Street,
		&address.House,
		&address.Apartment,
		&address.ZipCode,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &address, nil
}
