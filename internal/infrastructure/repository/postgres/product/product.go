package product

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

const (
	queryGetAllProducts = `
		SELECT id, seller_id, name, preview_image_url, description, 
				status, price, quantity, updated_at, rating, reviews_count 
			FROM product WHERE status = 'approved'
	`
	queryGetProductByID = `
		SELECT id, seller_id, name, preview_image_url, description, 
				status, price, quantity, updated_at, rating, reviews_count 
			FROM product WHERE id = $1
	`
	queryGetProductsByCategory = `
        SELECT p.id, p.seller_id, p.name, p.preview_image_url, p.description, 
                p.status, p.price, p.quantity, p.updated_at, p.rating, p.reviews_count 
			FROM product p
			JOIN product_category pc ON p.id = pc.product_id
			WHERE pc.category_id = $1 AND p.status = 'approved'
    `

	queryGetAllCategories =`
		SELECT id, name FROM category
	`
)

type ProductRepository struct {
	DB  *sql.DB
	log *logrus.Logger
}

// создание репозитория с заполнением данными
func NewProductRepository(db *sql.DB, log *logrus.Logger) *ProductRepository {
	return &ProductRepository{
		DB:  db,
		log: log,
	}
}

// получение основной информации всех товаров
func (p *ProductRepository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	productsList := []*models.Product{}

	rows, err := p.DB.QueryContext(ctx, queryGetAllProducts)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productsList, nil
}

// получение товара по id
func (p *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product := &models.Product{}
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
		)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepository) GetProductCoverPath(ctx context.Context, id uuid.UUID) ([]byte, error) {
	storagePath := models.GetProductCoverPath(id)

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cover image not found")
	}

	return os.ReadFile(storagePath)
}

func (p *ProductRepository) GetProductsByCategory(ctx context.Context, id uuid.UUID)([]*models.Product, error){
	productsList := []*models.Product{}

	rows, err := p.DB.QueryContext(ctx, queryGetProductsByCategory, id)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		productsList = append(productsList, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productsList, nil
}

func (p *ProductRepository) GetAllCategories(ctx context.Context)([]*models.Category, error){
	categoriesList := []*models.Category{}

	rows, err := p.DB.QueryContext(ctx, queryGetAllCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		category := &models.Category{}
		err = rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}
		categoriesList = append(categoriesList, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categoriesList, nil
}