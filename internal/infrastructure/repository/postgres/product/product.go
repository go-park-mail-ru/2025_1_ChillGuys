package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

const (
	queryGetAllProducts = `
		SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price
		FROM bazaar.product p 
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
		WHERE p.status = 'approved'
	`
	queryGetProductByID = `
		SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price
		FROM bazaar.product p
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
		WHERE p.id = $1
	`
	queryGetProductsByCategory = `
        SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
                p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count 
			FROM bazaar.product p
			JOIN bazaar.product_category pc ON p.id = pc.product_id
			WHERE pc.category_id = $1 AND p.status = 'approved'
    `
)

type ProductRepository struct {
	DB  *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		DB:  db,
	}
}

// получение основной информации всех товаров
func (p *ProductRepository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	const op = "ProductRepository.GetAllProducts"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.Product{}

	rows, err := p.DB.QueryContext(ctx, queryGetAllProducts)
	if err != nil {
		logger.WithError(err).Error("query all products")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var priceDiscount sql.NullFloat64
		product := &models.Product{}
		err = rows.Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.PreviewImageURL,
			&product.Description,
			&product.Status,
			&product.Price,
			&product.Quantity,
			&product.UpdatedAt,
			&product.Rating,
			&product.ReviewsCount,
			&priceDiscount,
		)
		if err != nil {
			logger.WithError(err).Error("scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		product.PriceDiscount = priceDiscount.Float64
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
        return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}

// получение товара по id
func (p *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	const op = "ProductRepository.GetProductByID"
    logger := logctx.GetLogger(ctx).WithField("op", op)
	
	product := &models.Product{}
	var priceDiscount sql.NullFloat64
	err := p.DB.QueryRowContext(ctx, queryGetProductByID, id).
		Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.PreviewImageURL,
			&product.Description,
			&product.Status,
			&product.Price,
			&product.Quantity,
			&product.UpdatedAt,
			&product.Rating,
			&product.ReviewsCount,
			&priceDiscount,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
            logger.Warn("product not found by ID")
            return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
        }
        logger.WithError(err).Error("failed to get product by ID")
        return nil, fmt.Errorf("%s: %w", op, err)
	}
	product.PriceDiscount = priceDiscount.Float64

	return product, nil
}

func (p *ProductRepository) GetProductsByCategory(ctx context.Context, id uuid.UUID) ([]*models.Product, error) {
	const op = "ProductRepository.GetProductsByCategory"
    logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.Product{}

	rows, err := p.DB.QueryContext(ctx, queryGetProductsByCategory, id)
	if err != nil {
		logger.WithError(err).Error("query products by category")
        return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		product := &models.Product{}
		err = rows.Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.PreviewImageURL,
			&product.Description,
			&product.Status,
			&product.Price,
			&product.Quantity,
			&product.UpdatedAt,
			&product.Rating,
			&product.ReviewsCount,
		)
		if err != nil {
			logger.WithError(err).Error("scan product row")
            return nil, fmt.Errorf("%s: %w", op, err)
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
        return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsList, nil
}