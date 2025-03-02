package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

type ProductRepo struct{
	Storage map[int]*models.Product
	Order []int 
	Mu sync.RWMutex
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		Storage: make(map[int]*models.Product),
		Order: make([]int, 0),
		Mu: sync.RWMutex{},
	}
}

func (p *ProductRepo) GetAllProducts(ctx context.Context) (models.Products, error) { //nolint:unparam
	productList := make(models.Products, 0, len(p.Storage))
	p.Mu.RLock()
	defer p.Mu.RUnlock()
	for i := range p.Order {
		productList = append(productList, p.Storage[i])
	}

	return productList, nil
}

func (p *ProductRepo) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
    product, exists := p.Storage[id]
    if !exists {
        return nil, fmt.Errorf("product with ID %d not found", id)
    }
    return product, nil
}