package product

import (
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
	"github.com/google/uuid"
)

type Response struct {
	Code    int
	Message string
	Data    any
	Meta    *response.Pagination
}

var (
	SuccessCreateResponse = Response{Code: 200, Message: "Product created successfully"}
	SuccessListResponse   = Response{Code: 200, Message: "Products fetched successfully"}
)

type ProductResponse struct {
	UUID          uuid.UUID `json:"productId"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"imageUrl"`
	Stock         int       `json:"stock"`
	Condition     Condition `json:"condition"`
	Tags          []string  `json:"tags"`
	IsPurchasable bool      `json:"isPurchasable"`
	Price         int       `json:"price"`
	PurchaseCount int       `json:"purchaseCount"`
}

func CreateProductResponse(product Product) ProductResponse {
	return ProductResponse{
		UUID:          product.UUID,
		Name:          product.Name,
		ImageURL:      product.ImageURL,
		Stock:         product.Stock,
		Condition:     product.Condition,
		Tags:          product.Tags,
		IsPurchasable: product.IsPurchasable,
		Price:         product.Price,
		PurchaseCount: product.PurchaseCount,
	}
}
