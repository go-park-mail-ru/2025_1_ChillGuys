package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	MediaFolder          = "./media"
	CoverName            = "cover.jpeg"
	DefaultPathCoverName = "product-default"
)

// Константы статусов товара
const (
	ProductPending  ProductStatus = "pending"  // Ожидает
	ProductRejected ProductStatus = "rejected" // Отказано
	ProductApproved ProductStatus = "approved" // Одобрено
)

// ProductStatus представляет статус товара
type ProductStatus int

type Product struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	SellerID        uuid.UUID     `json:"seller_id" db:"seller_id"`
	Name            string        `json:"name" db:"name"`
	PreviewImageURL string        `json:"preview_image_url,omitempty" db:"preview_image_url"`
	Description     string        `json:"description,omitempty" db:"description"`
	Status          ProductStatus `json:"status" db:"status"`
	Price           float64       `json:"price" db:"price"`
	Quantity        uint          `json:"quantity" db:"quantity"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
	Rating          uint          `json:"rating,omitempty" db:"rating"`
	ReviewsCount    uint          `json:"reviews_count" db:"reviews_count"`
}

type ProductDiscount struct {
	DiscountedPrice   float64   `db:"discounted_price"`
	DiscountEndDate   time.Time `db:"end_date"`
	DiscountStartDate time.Time `db:"start_date"`
}

// String возвращает строковое представление статуса
func (s ProductStatus) String() string {
	return [...]string{
		"pending",
		"rejected",
		"approved",
	}[s]
}

// ParseProductStatus преобразует строку в ProductStatus
func ParseProductStatus(status string) (ProductStatus, error) {
	switch status {
	case "pending":
		return ProductPending, nil
	case "rejected":
		return ProductRejected, nil
	case "approved":
		return ProductApproved, nil
	default:
		return ProductPending, fmt.Errorf("unknown product status: %s", status)
	}
}

// Scan реализует интерфейс sql.Scanner для чтения из БД
func (s *ProductStatus) Scan(value interface{}) error {
    if value == nil {
        *s = ProductPending
        return nil
    }

    var statusStr string
    
    // Обрабатываем разные типы, которые могут прийти из БД
    switch v := value.(type) {
    case string:
        statusStr = v
    case []byte:
        statusStr = string(v)
    default:
        return fmt.Errorf("failed to scan ProductStatus: unsupported type %T", value)
    }

    status, err := ParseProductStatus(statusStr)
    if err != nil {
        return err
    }
    *s = status
    return nil
}

// Value реализует интерфейс driver.Valuer для записи в БД
func (s ProductStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (p Product) MarshalJSON() ([]byte, error) {
    type Alias Product
    return json.Marshal(&struct {
        Status string `json:"status"`
        *Alias
    }{
        Status: p.Status.String(),
        Alias:  (*Alias)(&p),
    })
}