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

//функция заполнения тестовыми данными
func (r *ProductRepo) populateMockData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 1; i <= 15; i++ {
		product := models.Product{
			ID:           i,
			Name:         "Product " + string(rune('0'+i)), 
			Description:  "Description " + string(rune('0'+i)),
			Count:        uint(i * 10),
			Price:        uint(i * 100),
			ReviewsCount: uint(i * 5),
			Rating:       4.0 + float64(i%10)*0.1, // Рейтинг от 4.0 до 4.9
		}

		r.storage[product.ID] = &product
		r.order = append(r.order, product.ID)
	}
}

//создание репозитория с заполнением данными
func NewProductRepo() *ProductRepo {
	repo := &ProductRepo{
		storage: make(map[int]*models.Product),
		order: make([]int, 0),
		mu: sync.RWMutex{},
	}

	repo.populateMockData()

	return repo
}

//получение основной информации всех товаров
func (p *ProductRepo) GetAllProducts(ctx context.Context) ([]*models.Product, error) { //nolint:unparam
	productList := make([]*models.Product, 0, len(p.storage))
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, id := range p.order {
        productList = append(productList, p.storage[id])
    }

	return productList, nil
}

//получение товара по id
func (p *ProductRepo) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
    product, exists := p.storage[id]
    if !exists {
        return nil, fmt.Errorf("product with ID %d not found", id)
    }
	
    return product, nil
}