package models

import (
	"encoding/json"
	"fmt"
	"github.com/guregu/null"
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

type OrderPreviewProductDTO struct {
	ProductImageURL null.String
	ProductQuantity uint
}
