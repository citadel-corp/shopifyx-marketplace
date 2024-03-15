package product

import (
	bankaccount "github.com/citadel-corp/shopifyx-marketplace/internal/bank_account"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
	"github.com/google/uuid"
)

type Response struct {
	Code    int
	Message string
	Data    any
	Meta    *response.Pagination
	Error   error
}

var (
	SuccessCreateResponse      = Response{Code: 200, Message: "Product created successfully"}
	SuccessListResponse        = Response{Code: 200, Message: "ok"}
	SuccessPatchResponse       = Response{Code: 200, Message: "Product patched successfully"}
	SuccessGetResponse         = Response{Code: 200, Message: "ok"}
	SuccessPurchaseResponse    = Response{Code: 200, Message: "Product purchased successfully"}
	SuccessUpdateStockResponse = Response{Code: 200, Message: "Stock updated successfully"}
	SuccessDeleteResponse      = Response{Code: 200, Message: "Product deleted successfully"}
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

type SellerResponse struct {
	Name             string                            `json:"name"`
	ProductSoldTotal int                               `json:"productSoldTotal"`
	BankAccounts     []bankaccount.BankAccountResponse `json:"bankAccounts"`
}

func CreateSellerResponse(user *user.User, bankAccounts []*bankaccount.BankAccount) SellerResponse {
	if user == nil {
		return SellerResponse{}
	}

	accts := make([]bankaccount.BankAccountResponse, len(bankAccounts))
	for i, bankAccount := range bankAccounts {
		if bankAccount != nil {
			accts[i] = bankaccount.BankAccountResponse{
				BankAccountID:     bankAccount.UUID.String(),
				BankName:          bankAccount.BankName,
				BankAccountName:   bankAccount.BankAccountName,
				BankAccountNumber: bankAccount.BankAccountNumber,
			}
		}
	}

	return SellerResponse{
		Name:             user.Name,
		ProductSoldTotal: user.ProductSoldTotal,
		BankAccounts:     accts,
	}
}

type ProductDetailResponse struct {
	Product ProductResponse `json:"product"`
	Seller  SellerResponse  `json:"seller"`
}
