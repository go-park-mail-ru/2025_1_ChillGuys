package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	queryGetAllProducts = "SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product WHERE status = 'approved'"
	queryGetProductByID = "SELECT id, seller_id, name, preview_image_url, description, status, price, quantity, updated_at, rating, reviews_count FROM product WHERE id = $1"
)

//go:generate mockgen -source=product.go -destination=./mocks/product_repository_mock.go -package=mocks IProductRepository
type IProductRepository interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductCoverPath(ctx context.Context, id uuid.UUID) ([]byte, error)
}

type ProductRepository struct {
	DB *sql.DB
	log *logrus.Logger
}

//создание репозитория с заполнением данными
func NewProductRepository(db *sql.DB, log *logrus.Logger) *ProductRepository {
	return &ProductRepository{
		DB: db,
		log: log,
	}
}

//получение основной информации всех товаров
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

//получение товара по id
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

func (p *ProductRepository) GetProductCoverPath(ctx context.Context, id uuid.UUID) ([]byte, error){
	storagePath := models.GetProductCoverPath(id)

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cover image not found")
	}

	return os.ReadFile(storagePath)
}