package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

type ProductRepo struct{
	storage map[int]*models.Product
	order []int 
	mu sync.RWMutex
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		storage: make(map[int]*models.Product),
		order: make([]int, 0),
		mu: sync.RWMutex{},
	}
}

func (p *ProductRepo) GetAllProducts(ctx context.Context) ([]*models.Product, error) { //nolint:unparam
	productList := make([]*models.Product, 0, len(p.storage))
	p.mu.RLock()
	defer p.mu.RUnlock()

	for i := range p.order {
		productList = append(productList, p.storage[i])
	}

	return productList, nil
}

func (p *ProductRepo) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
    product, exists := p.storage[id]
    if !exists {
        return nil, fmt.Errorf("product with ID %d not found", id)
    }
	
    return product, nil
}