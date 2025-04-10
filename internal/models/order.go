package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

type OrderStatus int

const (
	Placed                    OrderStatus = iota // Оформлен
	AwaitingConfirmation                         // Ожидает подтверждения
	BeingPrepared                                // Готовится
	Shipped                                      // Отправлен
	InTransit                                    // В пути
	DeliveredToPickupPoint                       // Доставлен в пункт самовывоза
	Delivered                                    // Доставлен
	Canceled                                     // Отменен
	AwaitingPayment                              // Ожидает оплаты
	Paid                                         // Оплачено (опечатка в оригинале: должно быть Paid)
	PaymentFailed                                // Платеж не удался
	ReturnRequested                              // Возврат запрашивается
	ReturnProcessed                              // Возврат обработан
	ReturnInitiated                              // Возврат инициирован
	ReturnCompleted                              // Возврат завершен
	CanceledByUser                               // Отменен пользователем
	CanceledBySeller                             // Отменен продавцом
	CanceledDueToPaymentError                    // Отменен из-за ошибки платежа
)

func (s OrderStatus) String() string {
	return [...]string{
		"placed",
		"awaiting_confirmation",
		"being_prepared",
		"shipped",
		"in_transit",
		"delivered_to_pickup_point",
		"delivered",
		"canceled",
		"awaiting_payment",
		"paid",
		"payment_failed",
		"return_requested",
		"return_processed",
		"return_initiated",
		"return_completed",
		"canceled_by_user",
		"canceled_by_seller",
		"canceled_due_to_payment_error",
	}[s]
}

func ParseOrderStatus(s string) (OrderStatus, error) {
	statuses := [...]string{
		"pending",
		"placed",
		"awaiting_confirmation",
		"being_prepared",
		"shipped",
		"in_transit",
		"delivered_to_pickup_point",
		"delivered",
		"canceled",
		"awaiting_payment",
		"paid",
		"payment_failed",
		"return_requested",
		"return_processed",
		"return_initiated",
		"return_completed",
		"canceled_by_user",
		"canceled_by_seller",
		"canceled_due_to_payment_error",
	}

	for i, val := range statuses {
		if s == val {
			return OrderStatus(i), nil
		}
	}

	return 0, fmt.Errorf("unknown order status: %s", s)
}

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type OrderDB struct {
	ID         uuid.UUID `db:"id"`
	UserID     string    `db:"user_id"`
	Status     string    `db:"status"`
	TotalPrice float64   `db:"total_price"`
	AddressID  uuid.UUID `db:"address_id"`
}

type OrderItemDB struct {
	ID        uuid.UUID `db:"id"`
	OrderID   uuid.UUID `db:"order_id"`
	ProductID uuid.UUID `db:"product_id"`
	Price     float64   `db:"product_price"`
	Quantity  uint      `db:"quantity"`
}

type OrderPreview struct {
	ID                 uuid.UUID             `json:"id"`
	Status             OrderStatus           `json:"status"`
	TotalPrice         float64               `json:"total_price"`
	TotalDiscountPrice float64               `json:"total_discount_price"`
	Products           []OrderPreviewProduct `json:"products"`
	Address            AddressDB             `json:"address"`
	ExpectedDeliveryAt *time.Time            `json:"expected_delivery_at"`
	ActualDeliveryAt   *time.Time            `json:"actual_delivery_at"`
	CreatedAt          *time.Time            `json:"created_at,omitempty"`
}

type OrderPreviewProduct struct {
	ProductImageURL null.String `json:"product_image_url"`
	ProductQuantity uint        `json:"product_quantity"`
}
