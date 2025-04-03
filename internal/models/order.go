package models

import "github.com/google/uuid"

type OrderStatus int

const (
	Pending                   OrderStatus = iota // Ожидает
	Placed                                       // Оформлен
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
	}[s]
}

type Order struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Status             OrderStatus
	TotalPrice         int
	TotalPriceDiscount int
	AddressID          uuid.UUID
	Items              []CreateOrderItemDTO
}

type OrderDB struct {
	ID         uuid.UUID `db:"id"`
	UserID     string    `db:"user_id"`
	Status     string    `db:"status"`
	TotalPrice int       `db:"total_price"`
	AddressID  uuid.UUID `db:"address_id"`
}

type OrderItemDB struct {
	ID        uuid.UUID `db:"id"`
	OrderID   uuid.UUID `db:"order_id"`
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
}

type CreateOrderDTO struct {
	UserID    uuid.UUID
	Items     []CreateOrderItemDTO `json:"items"`
	AddressID uuid.UUID            `json:"address_id"`
}

type CreateOrderItemDTO struct {
	ID        uuid.UUID
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}
