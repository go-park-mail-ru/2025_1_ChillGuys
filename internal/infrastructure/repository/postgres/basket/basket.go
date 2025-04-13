package basket

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

const(
	queryGetBasketByUserID = `SELECT id FROM bazaar.basket WHERE user_id = $1 LIMIT 1`

	queryGetInfoBasket = `
		SELECT id, user_id, total_price, total_price_discount 
			FROM bazaar.basket WHERE id = $1
	`

	queryAddProductInBasket = `
		INSERT INTO bazaar.basket_item (id, basket_id, product_id, quantity)
			VALUES ($1, $2, $3, 1)
			ON CONFLICT (basket_id, product_id) 
			DO UPDATE SET 
				quantity = basket_item.quantity + 1
			RETURNING id, basket_id, product_id, quantity, updated_at
	`

	queryGetProductsInBasket = `
		SELECT 
			bi.id, 
			bi.basket_id, 
			bi.product_id, 
			bi.quantity, 
			bi.updated_at,
			p.name,
			p.price,
			p.preview_image_url,
			d.discounted_price
		FROM 
			bazaar.basket_item bi
		JOIN 
			bazaar.product p ON bi.product_id = p.id
		LEFT JOIN LATERAL (
			SELECT 
				discounted_price
			FROM 
				bazaar.discount
			WHERE 
				product_id = bi.product_id
				AND now() BETWEEN start_date AND end_date
			ORDER BY 
				start_date DESC
			LIMIT 1
		) d ON true
		WHERE 
			bi.basket_id = $1
			AND p.status = 'approved'
    `

	queryDelProductFromBasket = `
		DELETE FROM bazaar.basket_item
		WHERE basket_id = $1 AND product_id = $2
		RETURNING id
	`

	queryUpdateProductQuantity = `
		UPDATE bazaar.basket_item
		SET quantity = $1
		WHERE basket_id = $2 AND product_id = $3
		RETURNING id, basket_id, product_id, quantity, updated_at
	`

	queryGetQuantityProduct = `SELECT quantity FROM bazaar.product WHERE id = $1`

	queryClearBasket = `
		DELETE FROM bazaar.basket_item
		WHERE basket_id = $1
	`
)

type BasketRepository struct{
	DB  *sql.DB
}

func NewBasketRepository(db *sql.DB) *BasketRepository {
	return &BasketRepository{
		DB:  db,
	}
}

func (r *BasketRepository) getBasket(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	const op = "BasketRepository.getBasket"
    logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)
	
	var basketID uuid.UUID

	err := r.DB.QueryRowContext(ctx, queryGetBasketByUserID, userID).Scan(&basketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("basket not found")
            return uuid.Nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("get basket")
        return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return basketID, nil
}

func (r*BasketRepository) getQuantityProduct(ctx context.Context, productID uuid.UUID) (uint, error) {
	const op = "BasketRepository.getQuantityProduct"
    logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_id", productID)

	var quantity uint

	err := r.DB.QueryRowContext(ctx, queryGetQuantityProduct, productID).Scan(&quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("basket not found")
            return 0, errs.NewNotFoundError(op)
        }
        logger.WithError(err).Error("get basket")
        return 0, fmt.Errorf("%s: %w", op, err)
	}

	return quantity, nil
}

func (r *BasketRepository) Get(ctx context.Context, userID uuid.UUID) ([]*models.BasketItem, error) {
	const op = "BasketRepository.Get"
    logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	basketID, err := r.getBasket(ctx, userID)
	if err != nil {
		logger.WithError(err).Error("get basket ID")
        return nil, fmt.Errorf("%s: %w", op, err)
	}

	productsList := []*models.BasketItem{}

	rows, err := r.DB.QueryContext(ctx, queryGetProductsInBasket, basketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("no products in basket")
            return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("query basket products")
        return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.BasketItem{}
		var priceDiscount sql.NullFloat64
		err = rows.Scan(
			&item.ID,
			&item.BasketID,
			&item.ProductID,
			&item.Quantity,
			&item.UpdatedAt,
			&item.ProductName,
			&item.Price,
			&item.ProductImage,
			&priceDiscount,
		)
		if err != nil {
			logger.WithError(err).Error("scan basket item")
            return nil, fmt.Errorf("%s: %w", op, err)
		}
		item.PriceDiscount = priceDiscount.Float64
		productsList = append(productsList, item)
	}

	if err = rows.Err(); err != nil {
        logger.WithError(err).Error("rows iteration error")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return productsList, nil
}

func (r *BasketRepository) Add(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*models.BasketItem, error){
	const op = "BasketRepository.Add"
    logger := logctx.GetLogger(ctx).WithField("op", op).
        WithField("user_id", userID).
        WithField("product_id", productID)

	basketID, err := r.getBasket(ctx, userID)
    if err != nil {
        logger.WithError(err).Error("get basket ID")
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	item := &models.BasketItem{}
	newItemID := uuid.New()

	err = r.DB.QueryRowContext(ctx, queryAddProductInBasket, newItemID, basketID, productID).Scan(
		&item.ID,
		&item.BasketID,
		&item.ProductID,
		&item.Quantity,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("add product to basket - not found")
            return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("add product to basket")
        return nil, fmt.Errorf("%s: %w", op, err)
	}

	return item, nil
}

func (r *BasketRepository) Delete(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	const op = "BasketRepository.Delete"
    logger := logctx.GetLogger(ctx).WithField("op", op).
        WithField("user_id", userID).
        WithField("product_id", productID)

	basketID, err := r.getBasket(ctx, userID)
	if err != nil {
        logger.WithError(err).Error("get basket ID")
        return fmt.Errorf("%s: %w", op, err)
    }

	var deletedID uuid.UUID
	err = r.DB.QueryRowContext(ctx, queryDelProductFromBasket, basketID, productID).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("product not found in basket")
            return fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("delete product from basket")
        return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *BasketRepository) UpdateQuantity(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.BasketItem, int, error) {
	const op = "BasketRepository.UpdateQuantity"
    logger := logctx.GetLogger(ctx).WithField("op", op).
        WithField("user_id", userID).
        WithField("product_id", productID).
        WithField("quantity", quantity)

	quantityProduct, err := r.getQuantityProduct(ctx, productID)
	if err != nil {
		logger.WithError(err).Error("failed to get basket ID")
        return nil, -1, fmt.Errorf("%s: %w", op, err)
	}

	if uint(quantity) > quantityProduct {
        logger.Warn("requested quantity exceeds available")
        return nil, -1, errs.NewBusinessLogicError("requested quantity exceeds available stock")
    }

	basketID, err := r.getBasket(ctx, userID)
	if err != nil {
        logger.WithError(err).Error("failed to get basket ID")
        return nil, -1, fmt.Errorf("%s: %w", op, err)
    }

	item := &models.BasketItem{}
	err = r.DB.QueryRowContext(ctx, queryUpdateProductQuantity, quantity, basketID, productID).Scan(
		&item.ID,
		&item.BasketID,
		&item.ProductID,
		&item.Quantity,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("product not found in basket")
            return nil, -1, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("update product quantity")
        return nil, -1, fmt.Errorf("%s: %w", op, err)
	}

	return item, int(quantityProduct)-quantity, nil
}

func (r *BasketRepository) Clear(ctx context.Context, userID uuid.UUID) error {
	const op = "BasketRepository.Clear"
    logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	basketID, err := r.getBasket(ctx, userID)
	if err != nil {
        logger.WithError(err).Error("get basket ID")
        return fmt.Errorf("%s: %w", op, err)
    }

	_, err = r.DB.ExecContext(ctx, queryClearBasket, basketID)
	if err != nil {
        logger.WithError(err).Error("clear basket")
        return fmt.Errorf("%s: %w", op, err)
    }

	return nil
}