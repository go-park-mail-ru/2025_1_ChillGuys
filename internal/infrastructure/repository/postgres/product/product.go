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
        ORDER BY p.id
		LIMIT 20 OFFSET $1
	`
	queryGetProductByID = `
		SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
				p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count,
				d.discounted_price, s.id, s.title, s.description
		FROM bazaar.product p
		LEFT JOIN bazaar.discount d ON p.id = d.product_id
        LEFT JOIN bazaar.seller s ON s.id = p.seller_id
		WHERE p.id = $1
	`

	queryGetProductsByCategoryWithFilterAndSort = `
        SELECT 
            p.id, 
            p.seller_id, 
            p.name, 
            p.preview_image_url, 
            p.description, 
            p.status, 
            p.price, 
            p.quantity, 
            p.updated_at, 
            p.rating, 
            p.reviews_count
        FROM 
            bazaar.product p
        JOIN 
            bazaar.product_subcategory pc ON p.id = pc.product_id
        WHERE 
            pc.subcategory_id = $1 
            AND p.status = 'approved'
            AND ($3 = 0 OR p.price > $3)
            AND ($4 = 0 OR p.price < $4)
            AND ($5 = 0::FLOAT OR p.rating > $5::FLOAT)
        ORDER BY %s
        LIMIT 20 OFFSET $2
    `

	queryAddProduct = `
		INSERT INTO bazaar.product (
			id, seller_id, name, preview_image_url, 
			description, status, price, quantity, rating, reviews_count
		) VALUES ($1, $2, $3, $4, $5, $10, $6, $7, $8, $9)
		RETURNING id
	`

	queryAddDiscount = `
        INSERT INTO bazaar.discount (
            id, product_id, discounted_price, start_date, end_date
        ) VALUES ($1, $2, $3, now(), now() + interval '30 days')
    `

	queryAddProductCategory = `
        INSERT INTO bazaar.product_subcategory (id, product_id, subcategory_id)
        VALUES ($1, $2, $3)
    `
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

// получение основной информации всех товаров
func (p *ProductRepository) GetAllProducts(ctx context.Context, offset int) ([]*models.Product, error) {
	const op = "ProductRepository.GetAllProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	productsList := []*models.Product{}

	rows, err := p.DB.QueryContext(ctx, queryGetAllProducts, offset)
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
	logger.Print(id)
	product := &models.Product{}
	var seller models.Seller
	var sellerID uuid.NullUUID
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
			&sellerID,
			&seller.Title,
			&seller.Description,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("product not found by ID")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("failed to get product by ID")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if sellerID.Valid {
		seller.ID = sellerID.UUID
		product.Seller = &seller
	}
	product.PriceDiscount = priceDiscount.Float64

	return product, nil
}

func (p *ProductRepository) GetProductsByCategory(
	ctx context.Context,
	id uuid.UUID,
	offset int,
	minPrice, maxPrice float64,
	minRating float32,
	sortOption models.SortOption,
) ([]*models.Product, error) {
	const op = "ProductRepository.GetProductsByCategoryWithFilterAndSort"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var orderByClause string
	switch sortOption {
	case models.SortByPriceAsc:
		orderByClause = "p.price ASC"
	case models.SortByPriceDesc:
		orderByClause = "p.price DESC"
	case models.SortByRatingAsc:
		orderByClause = "p.rating ASC"
	case models.SortByRatingDesc:
		orderByClause = "p.rating DESC"
	default:
		orderByClause = "p.updated_at DESC" // дефолтная сортировка
	}

	query := fmt.Sprintf(queryGetProductsByCategoryWithFilterAndSort, orderByClause)

	productsList := []*models.Product{}

	rows, err := p.DB.QueryContext(
		ctx,
		query,
		id,
		offset,
		minPrice,
		maxPrice,
		minRating,
	)

	if err != nil {
		logger.WithError(err).Error("query products by category with filter and sort")
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

func (p *ProductRepository) AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error) {
	const op = "ProductRepository.AddProduct"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		logger.WithError(err).Error("begin transaction")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	// Генерируем UUID для продукта
	product.ID = uuid.New()

	// Вставляем продукт
	product.Status = models.ProductApproved
	_, err = tx.ExecContext(ctx, queryAddProduct,
		product.ID,
		product.SellerID,
		product.Name,
		product.PreviewImageURL,
		product.Description,
		product.Price,
		product.Quantity,
		product.Rating,
		product.ReviewsCount,
		product.Status,
	)
	if err != nil {
		logger.WithError(err).Error("insert product")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Генерируем UUID для скидки и добавляем скидку
	discountID := uuid.New()
	_, err = tx.ExecContext(ctx, queryAddDiscount,
		discountID,
		product.ID,
		product.PriceDiscount,
	)
	if err != nil {
		logger.WithError(err).Error("insert discount")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	productCategoryID := uuid.New()
	// Добавляем категорию продукта
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
