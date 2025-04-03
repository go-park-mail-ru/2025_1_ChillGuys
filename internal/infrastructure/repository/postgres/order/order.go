package order

import (
	"context"
	"database/sql"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	queryCreateOrder  = `INSERT INTO "order" (id, user_id, status, total_price, total_price_discount, address_id) VALUES ($1, $2, $3, $4, $5, $6)`
	queryAddOrderItem = `INSERT INTO "order_item" (id, order_id, product_id, quantity) VALUES ($1, $2, $3, $4)`
)

type IOrderRepository interface {
	CreateOrder(context.Context, models.Order) error
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

func (r *OrderRepository) CreateOrder(ctx context.Context, in models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryCreateOrder,
		in.ID, in.UserID, in.Status.String(), in.TotalPrice, in.TotalPriceDiscount, in.AddressID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range in.Items {
		_, err = tx.ExecContext(ctx, queryAddOrderItem,
			item.ID, in.ID, item.ProductID, item.Quantity,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
