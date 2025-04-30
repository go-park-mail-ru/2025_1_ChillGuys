package seller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

const (
	queryAddProduct = `
		INSERT INTO bazaar.product (
			id, seller_id, name, 
			description, status, price, quantity, rating, reviews_count
		) VALUES ($1, $2, $3, $4, $5, $6, $7, 0, 0)
		RETURNING id
	`

	queryAddProductCategory = `
		INSERT INTO bazaar.product_subcategory (id, product_id, subcategory_id)
		VALUES ($1, $2, $3)
	`

	queryGetSellerProducts = `
		SELECT id, seller_id, name, preview_image_url, 
			description, status, price, quantity, rating, reviews_count
		FROM bazaar.product
		WHERE seller_id = $1
		LIMIT 20 OFFSET $2
	`

	queryCheckProductBelongs = `
		SELECT EXISTS(
			SELECT 1 FROM bazaar.product 
			WHERE id = $1 AND seller_id = $2
		)
	`

	queryUpdateProductImage = `
		UPDATE bazaar.product
		SET preview_image_url = $1
		WHERE id = $2
		RETURNING preview_image_url
	`
)

type SellerRepository struct {
	db *sql.DB
}

func NewSellerRepository(db *sql.DB) *SellerRepository {
	return &SellerRepository{db: db}
}

func (r *SellerRepository) AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error) {
	const op = "SellerRepository.AddProduct"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("begin transaction")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	// Генерируем UUID для продукта
	product.ID = uuid.New()

	// Вставляем продукт
	product.Status = models.ProductPending
	_, err = tx.ExecContext(ctx, queryAddProduct,
		product.ID,
		product.SellerID,
		product.Name,
		product.Description,
		product.Status,
		product.Price,
		product.Quantity,
	)
	if err != nil {
		logger.WithError(err).Error("insert product")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Добавляем категорию продукта
	productCategoryID := uuid.New()
	_, err = tx.ExecContext(ctx, queryAddProductCategory,
		productCategoryID,
		product.ID,
		categoryID,
	)
	if err != nil {
		logger.WithError(err).Error("insert product category")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		logger.WithError(err).Error("commit transaction")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return product, nil
}

func (r *SellerRepository) UploadProductImage(ctx context.Context, productID uuid.UUID, imageURL string) error {
	const op = "SellerRepository.UploadProductImage"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	// Обновляем URL изображения в БД
	_, err := r.db.ExecContext(ctx, queryUpdateProductImage, imageURL, productID)
	if err != nil {
		logger.WithError(err).Error("update product image URL")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *SellerRepository) GetSellerProducts(ctx context.Context, sellerID uuid.UUID, offset int) ([]*models.Product, error) {
	const op = "SellerRepository.GetSellerProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetSellerProducts, sellerID, offset)
	if err != nil {
		logger.WithError(err).Error("query seller products")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.PreviewImageURL,
			&product.Description,
			&product.Status,
			&product.Price,
			&product.Quantity,
			&product.Rating,
			&product.ReviewsCount,
		)
		if err != nil {
			logger.WithError(err).Error("scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

func (r *SellerRepository) CheckProductBelongs(ctx context.Context, productID, sellerID uuid.UUID) (bool, error) {
	const op = "SellerRepository.CheckProductBelongs"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var belongs bool
	err := r.db.QueryRowContext(ctx, queryCheckProductBelongs, productID, sellerID).Scan(&belongs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		logger.WithError(err).Error("check product belongs")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return belongs, nil
}