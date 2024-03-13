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
