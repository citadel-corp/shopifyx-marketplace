package product

import (
	"math"

	validation "github.com/itgelo/ozzo-validation/v4"
	"github.com/itgelo/ozzo-validation/v4/is"
)

type CreateProductPayload struct {
	Name          string   `json:"name"`
	Price         int      `json:"price"`
	ImageURL      string   `json:"imageUrl"`
	Stock         int      `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchaseable"`
	UserID        uint64   `json:"-"`
}

func (p CreateProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required.Error(ErrorRequiredField.Message), validation.Length(5, 60)),
		validation.Field(&p.Price, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.ImageURL, validation.Required.Error(ErrorRequiredField.Message), is.URL),
		validation.Field(&p.Stock, validation.Required.Error(ErrorRequiredField.Message), validation.Min(0)),
		validation.Field(&p.Condition, validation.Required.Error(ErrorRequiredField.Message), validation.In(New.String(), Second.String())),
		validation.Field(&p.Tags, validation.Required.Error(ErrorRequiredField.Message), validation.Length(0, math.MaxInt)),
		validation.Field(&p.IsPurchasable, validation.Required.Error(ErrorRequiredField.Message)),
		validation.Field(&p.UserID, validation.Required.Error(ErrorUnauthorized.Message)),
	)
}

type ListProductPayload struct {
	UserOnly       bool
	UserID         uint64
	Tags           []string
	Condition      string
	ShowEmptyStock bool
	MinPrice       int
	MaxPrice       int
	Search         string
	Limit          int
	Offset         int
	SortBy         string
	OrderBy        string
}

func (p ListProductPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.UserID, validation.When(p.UserOnly, validation.Required.Error(ErrorUnauthorized.Message))),
		validation.Field(&p.Condition, validation.In(New.String(), Second.String())),
		validation.Field(&p.MinPrice, validation.When(p.MaxPrice != 0, validation.Max(p.MaxPrice))),
		validation.Field(&p.MaxPrice, validation.When(p.MinPrice != 0, validation.Min(p.MinPrice))),
		validation.Field(&p.SortBy, validation.In("price", "date")),
		validation.Field(&p.OrderBy, validation.In("asc", "dsc")),
		validation.Field(&p.Limit, validation.When(p.Offset != 0, validation.Required.Error(ErrorRequiredField.Message))),
		validation.Field(&p.Offset, validation.When(p.Limit != 0, validation.Required.Error(ErrorRequiredField.Message))),
	)
}
