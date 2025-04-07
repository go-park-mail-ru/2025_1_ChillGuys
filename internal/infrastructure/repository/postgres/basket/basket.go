package basket

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

//проверяем есть ли у юзера корзина если нет то создаем (по юзер айди)

//добвление товара в корзину - создание элемента в таблице basket_item
//quantity = 1

//отдача всех товаров в корзине т.е всех basket_item в которых basket_id соответствует юзер_id в basket_user


const(
	queryGetOrCreateBasket = `
		WITH 
		existing_basket AS (
			SELECT id FROM basket 
			WHERE user_id = $1
			LIMIT 1
		),
		new_basket AS (
			INSERT INTO basket (id, user_id, total_price, total_price_discount)
			SELECT $2, $1, 0, 0
			WHERE NOT EXISTS (SELECT 1 FROM existing_basket)
			RETURNING id
		)
		SELECT id FROM existing_basket
		UNION ALL
		SELECT id FROM new_basket
		LIMIT 1
	`

	queryGetInfoBasket = `
		SELECT id, user_id, total_price, total_price_discount 
			FROM basket WHERE id = $1
	`

	queryAddProductInBasket = `
		INSERT INTO basket_item (id, basket_id, product_id, quantity)
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
			basket_item bi
		JOIN 
			product p ON bi.product_id = p.id
		LEFT JOIN LATERAL (
			SELECT 
				discounted_price
			FROM 
				discount
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
		DELETE FROM basket_item
		WHERE basket_id = $1 AND product_id = $2
		RETURNING id
	`

	queryUpdateProductQuantity = `
		UPDATE basket_item
		SET quantity = $1
		WHERE basket_id = $2 AND product_id = $3
		RETURNING id, basket_id, product_id, quantity, updated_at
	`

	queryClearBasket = `
		DELETE FROM basket_item
		WHERE basket_id = $1
	`
)

type BasketRepository struct{
	DB  *sql.DB
	log *logrus.Logger
}

func NewBasketRepository(db *sql.DB, log *logrus.Logger) *BasketRepository {
	return &BasketRepository{
		DB:  db,
		log: log,
	}
}

func (r *BasketRepository) GetOrCreateBasket(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var basketID uuid.UUID
	newBasketID := uuid.New()

	err := r.DB.QueryRowContext(ctx, queryGetOrCreateBasket, userID, newBasketID).Scan(&basketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errs.ErrNotFound
		}
		return uuid.Nil, err
	}

	return basketID, nil
}

func (r *BasketRepository) GetProductsInBasket(ctx context.Context, userID uuid.UUID) (*dto.BasketResponse, error) {
	basketID, err := r.GetOrCreateBasket(ctx, userID)
	if err != nil {
		return nil, err
	}

	productsList := []*models.BasketItem{}

	rows, err := r.DB.QueryContext(ctx, queryGetProductsInBasket, basketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
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
			return nil, err
		}
		item.PriceDiscount = priceDiscount.Float64
		productsList = append(productsList, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	response := dto.ConvertToBasketResponse(productsList)

	return &response, nil
}

func (r *BasketRepository) AddProductInBasket(ctx context.Context, userID uuid.UUID, productID uuid.UUID) (*models.BasketItem, error){
	basketID, err := r.GetOrCreateBasket(ctx, userID)
    if err != nil {
        return nil, err
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
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *BasketRepository) DeleteProductFromBasket(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	basketID, err := r.GetOrCreateBasket(ctx, userID)
	if err != nil {
		return err
	}

	var deletedID uuid.UUID
	err = r.DB.QueryRowContext(ctx, queryDelProductFromBasket, basketID, productID).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}

	return nil
}

func (r *BasketRepository) UpdateProductQuantity(ctx context.Context, userID uuid.UUID, productID uuid.UUID, quantity int) (*models.BasketItem, error) {
	basketID, err := r.GetOrCreateBasket(ctx, userID)
	if err != nil {
		return nil, err
	}

	if quantity <= 0 {
		return nil, errs.ErrInvalidQuantity
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
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *BasketRepository) ClearBasket(ctx context.Context, userID uuid.UUID) error {
	basketID, err := r.GetOrCreateBasket(ctx, userID)
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, queryClearBasket, basketID)
	if err != nil {
		return err
	}

	return nil
}