package product

import (
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type Service struct {
	repository Repository
}

func (s *Service) Create(req CreateProductPayload) *Error {
	product := &Product{
		ImageURL:      req.ImageURL,
		Stock:         req.Stock,
		Condition:     ToCondition(req.Condition),
		Tags:          req.Tags,
		IsPurchasable: req.IsPurchasable,
		Price:         req.Price,
		User: user.User{
			ID: req.UserID,
		},
	}

	err := s.repository.Create(product)
	if err != nil {
		return &ErrorInternal
	}

	return nil
}
