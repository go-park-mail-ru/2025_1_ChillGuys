package dto

import "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"

type UpdateQuantityRequest struct {
	Quantity  int       `json:"quantity"`
}

type UpdateQuantityResponse struct {
	QuantityRemain 		int 				`json:"remain_quantity"`
	Item 				*models.BasketItem  `json:"item"`
}

func ConvertToQuantityResponse(item *models.BasketItem, quantity int) UpdateQuantityResponse {
	return UpdateQuantityResponse{
		QuantityRemain: quantity,
		Item: item,
	}
}

type BasketResponse struct {
    Total              int          		`json:"total"`
    TotalPrice         float64      		`json:"total_price"`
    TotalPriceDiscount float64      		`json:"total_price_discount"`
	Products           []models.BasketItem  `json:"products"`
}

func ConvertToBasketResponse(items []*models.BasketItem) BasketResponse{
	var totalPrice, totalPriceDiscount float64 
	productsList := make([]models.BasketItem, 0, len(items))
	for _, product := range items{
		productsList = append(productsList, *product)
		totalPrice = totalPrice + product.Price * float64(product.Quantity)
        price := product.Price
        if product.PriceDiscount > 0 {
            price = product.PriceDiscount
        }
        totalPriceDiscount += price * float64(product.Quantity)
	}

	return BasketResponse{
		Total: len(productsList),
		TotalPrice: totalPrice,
		TotalPriceDiscount : totalPriceDiscount,
		Products: productsList,
	}
}