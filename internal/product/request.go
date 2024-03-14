package product

import (
	"math"

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
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required.Error(ErrorRequiredField.Message), validation.Length(5, 60)),
		validation.Field(&p.Price, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.ImageURL, validation.Required.Error(ErrorRequiredField.Message), is.URL),
		validation.Field(&p.Stock, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.Condition, validation.Required.Error(ErrorRequiredField.Message), validation.In(Conditions...)),
		validation.Field(&p.Tags, validation.Required.Error(ErrorRequiredField.Message), validation.Length(0, math.MaxInt)),
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
		validation.Field(&p.OrderBy, validation.In("asc", "dsc")),
		validation.Field(&p.Limit, validation.When(p.Offset != 0, validation.Required.Error(ErrorRequiredField.Message))),
		validation.Field(&p.Offset, validation.When(p.Limit != 0, validation.NotNil.Error(ErrorRequiredField.Message))),
	)
}
