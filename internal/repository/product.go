package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	queryGetAllProducts = "SELECT id, name, description, count, price, reviews_count, rating FROM product"
	queryGetProductByID = "SELECT id, name, description, count, price, reviews_count, rating FROM product WHERE id = $1"
)

// type ProductRepository struct{
// 	storage map[int]*models.Product
// 	order []int
// 	mu sync.RWMutex
// }

type ProductRepository struct {
	DB *sql.DB
	log *logrus.Logger
}

// var products = []models.Product{
// 	{ID: 1, Name: "Смартфон Xiaomi Redmi Note 10", Description: "Смартфон с AMOLED-дисплеем и камерой 48 Мп", Count: 50, Price: 19999, ReviewsCount: 120, Rating: 4.5},
// 	{ID: 2, Name: "Ноутбук ASUS VivoBook 15", Description: "Ноутбук с процессором Intel Core i5 и SSD на 512 ГБ", Count: 30, Price: 54999, ReviewsCount: 80, Rating: 4.7},
// 	{ID: 3, Name: "Наушники Sony WH-1000XM4", Description: "Беспроводные наушники с шумоподавлением", Count: 25, Price: 29999, ReviewsCount: 200, Rating: 4.8},
// 	{ID: 4, Name: "Фитнес-браслет Xiaomi Mi Band 6", Description: "Фитнес-браслет с AMOLED-дисплеем и мониторингом сна", Count: 100, Price: 3999, ReviewsCount: 300, Rating: 4.6},
// 	{ID: 5, Name: "Пылесос Dyson V11", Description: "Беспроводной пылесос с мощным всасыванием", Count: 15, Price: 59999, ReviewsCount: 90, Rating: 4.9},
// 	{ID: 6, Name: "Кофемашина DeLonghi Magnifica", Description: "Автоматическая кофемашина для приготовления эспрессо", Count: 10, Price: 79999, ReviewsCount: 70, Rating: 4.7},
// 	{ID: 7, Name: "Электросамокат Xiaomi Mi Scooter 3", Description: "Электросамокат с запасом хода 30 км", Count: 40, Price: 29999, ReviewsCount: 150, Rating: 4.5},
// 	{ID: 8, Name: "Умная колонка Яндекс.Станция Мини", Description: "Умная колонка с голосовым помощником Алисой", Count: 60, Price: 7999, ReviewsCount: 250, Rating: 4.4},
// 	{ID: 9, Name: "Монитор Samsung Odyssey G5", Description: "Игровой монитор с разрешением 1440p и частотой 144 Гц", Count: 20, Price: 34999, ReviewsCount: 100, Rating: 4.6},
// 	{ID: 10, Name: "Электрочайник Bosch TWK 3A011", Description: "Электрочайник с мощностью 2400 Вт", Count: 50, Price: 1999, ReviewsCount: 180, Rating: 4.3},
// 	{ID: 11, Name: "Робот-пылесос iRobot Roomba 981", Description: "Робот-пылесос с навигацией по карте помещения", Count: 12, Price: 69999, ReviewsCount: 60, Rating: 4.8},
// 	{ID: 12, Name: "Фен Dyson Supersonic", Description: "Фен с технологией защиты волос от перегрева", Count: 18, Price: 49999, ReviewsCount: 130, Rating: 4.7},
// 	{ID: 13, Name: "Микроволновая печь LG MS-2042DB", Description: "Микроволновка с объемом 20 литров", Count: 35, Price: 8999, ReviewsCount: 110, Rating: 4.2},
// 	{ID: 14, Name: "Игровая консоль PlayStation 5", Description: "Игровая консоль нового поколения", Count: 5, Price: 79999, ReviewsCount: 300, Rating: 4.9},
// 	{ID: 15, Name: "Электронная книга PocketBook 740", Description: "Электронная книга с экраном E Ink Carta", Count: 25, Price: 19999, ReviewsCount: 90, Rating: 4.4},
// }

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
            &product.Name,
            &product.Description,
            &product.Count,
            &product.Price,
            &product.ReviewsCount,
            &product.Rating,
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
func (p *ProductRepository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	product := &models.Product{}
	err := p.DB.QueryRowContext(ctx, queryGetProductByID, id).
		Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Count,
			&product.Price,
			&product.ReviewsCount,
			&product.Rating,)
	
	if err != nil {
		return nil, err
	}
	
    return product, nil
}

func (p *ProductRepository) GetProductCoverPath(ctx context.Context, id int) ([]byte, error){
	storagePath := models.GetProductCoverPath(id)

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cover image not found")
	}

	return os.ReadFile(storagePath)
}