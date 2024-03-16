package product

import (
	"errors"

	"github.com/google/uuid"
	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/itgelo/ozzo-validation/v4/is"
)

type CreateProductPayload struct {
	Name          string    `json:"name"`
	Price         int       `json:"price"`
	ImageURL      string    `json:"imageUrl"`
	Stock         int       `json:"stock"`
	Condition     Condition `json:"condition"`
	Tags          []string  `json:"tags"`
	IsPurchasable bool      `json:"isPurchasable"`
	UserID        uint64    `json:"-"`
}

func (p CreateProductPayload) Validate() error {
	for i := range p.Tags {
		if len(p.Tags[i]) == 0 {
			return errors.New("tags must not be empty")
		}
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required.Error(ErrorRequiredField.Message), validation.Length(5, 60)),
		validation.Field(&p.Price, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.ImageURL, validation.Required.Error(ErrorRequiredField.Message), is.URL),
		validation.Field(&p.Stock, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.Condition, validation.Required.Error(ErrorRequiredField.Message), validation.In(Conditions...)),
		validation.Field(&p.Tags, validation.Required.Error(ErrorRequiredField.Message)),
		validation.Field(&p.IsPurchasable, validation.Required.Error(ErrorRequiredField.Message)),
		validation.Field(&p.UserID, validation.Required.Error(ErrorUnauthorized.Message)),
	)
}

type sortBy string
type productSortBy sortBy

var (
	SortByPrice productSortBy = "price"
	SortByDate  productSortBy = "date"
)

var productSortBys []interface{} = []interface{}{SortByPrice, SortByDate}

type ListProductPayload struct {
	UserOnly       bool `schema:"userOnly" binding:"omitempty"`
	UserID         uint64
	Tags           []string      `schema:"tags" binding:"omitempty"`
	Condition      Condition     `schema:"condition" binding:"omitempty"`
	ShowEmptyStock bool          `schema:"showEmptyStock" binding:"omitempty"`
	MinPrice       int           `schema:"minPrice" binding:"omitempty"`
	MaxPrice       int           `schema:"maxPrice" binding:"omitempty"`
	Search         string        `schema:"search" binding:"omitempty"`
	Limit          int           `schema:"limit" binding:"omitempty"`
	Offset         int           `schema:"offset" binding:"omitempty"`
	SortBy         productSortBy `schema:"sortBy" binding:"omitempty"`
	OrderBy        string        `schema:"orderBy" binding:"omitempty"`
}

func (p ListProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserID, validation.When(p.UserOnly, validation.Required.Error(ErrorUnauthorized.Message))),
		validation.Field(&p.Condition, validation.In(Conditions...)),
		validation.Field(&p.MinPrice, validation.When(p.MaxPrice != 0, validation.Max(p.MaxPrice))),
		validation.Field(&p.MaxPrice, validation.When(p.MinPrice != 0, validation.Min(p.MinPrice))),
		validation.Field(&p.SortBy, validation.In(productSortBys...)),
		validation.Field(&p.OrderBy, validation.In("asc", "desc")),
		validation.Field(&p.Limit, validation.When(p.Offset != 0, validation.Required.Error(ErrorRequiredField.Message))),
		validation.Field(&p.Offset, validation.When(p.Limit != 0, validation.NotNil.Error(ErrorRequiredField.Message))),
	)
}

type UpdateProductPayload struct {
	ProductUID    uuid.UUID `json:"-"`
	Name          string    `json:"name,omitempty"`
	Price         int       `json:"price,omitempty"`
	ImageURL      string    `json:"imageUrl,omitempty"`
	Condition     Condition `json:"condition,omitempty"`
	Tags          []string  `json:"tags,omitempty"`
	IsPurchasable bool      `json:"isPurchasable,omitempty"`
	UserID        uint64    `json:"-"`
}

func (p UpdateProductPayload) Validate() error {
	for i := range p.Tags {
		if len(p.Tags[i]) == 0 {
			return errors.New("tags must not be empty")
		}
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductUID, validation.Required.Error(ErrorRequiredField.Message)),
		validation.Field(&p.Name, validation.Required.Error(ErrorRequiredField.Message), validation.Length(5, 60)),
		validation.Field(&p.Price, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.ImageURL, validation.Required.Error(ErrorRequiredField.Message), is.URL),
		validation.Field(&p.Condition, validation.Required.Error(ErrorRequiredField.Message), validation.In(Conditions...)),
		validation.Field(&p.Tags, validation.Required.Error(ErrorRequiredField.Message)),
		validation.Field(&p.IsPurchasable, validation.Required.Error(ErrorRequiredField.Message)),
		validation.Field(&p.UserID, validation.Required.Error(ErrorUnauthorized.Message)),
	)
}

type GetProductPayload struct {
	ProductUID uuid.UUID `json:"-"`
}

func (p GetProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ProductUID, validation.Required.Error(ErrorRequiredField.Message)),
	)
}

type PurchaseProductPayload struct {
	ProductUID           uuid.UUID
	BankAccountID        uuid.UUID `json:"bankAccountId"`
	PaymentProofImageURL string    `json:"paymentProofImageUrl"`
	Quantity             int       `json:"quantity"`
	BuyerID              uint64
	SellerID             uint64
}

func (p PurchaseProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BankAccountID, validation.Required.Error(ErrorRequiredField.Message), is.UUID),
		validation.Field(&p.PaymentProofImageURL, validation.Required.Error(ErrorRequiredField.Message), is.URL),
		validation.Field(&p.Quantity, validation.Required.Error(ErrorRequiredField.Message), validation.Min(1)),
		validation.Field(&p.BuyerID, validation.Required.Error(ErrorUnauthorized.Message)),
	)
}

type UpdateStockPayload struct {
	ProductUID uuid.UUID
	Stock      int `json:"stock"`
	UserID     uint64
}

func (p UpdateStockPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Stock, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.UserID, validation.Required.Error(ErrorUnauthorized.Message)),
	)
}

type DeleteProductPayload struct {
	ProductUID uuid.UUID
	UserID     uint64
}

func (p DeleteProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserID, validation.Required.Error(ErrorUnauthorized.Message)),
	)
}
