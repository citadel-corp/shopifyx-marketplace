package product

import (
	"context"
	"log/slog"

	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type ProductService struct {
	repository Repository
}

type Service interface {
	Create(ctx context.Context, req CreateProductPayload) Response
}

func NewService(repository Repository) Service {
	return &ProductService{repository: repository}
}

func (s *ProductService) Create(ctx context.Context, req CreateProductPayload) Response {
	product := &Product{
		Name:          req.Name,
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

	err := s.repository.Create(ctx, product)
	if err != nil {
		slog.Error(err.Error())
		return ErrorInternal
	}

	return SuccessResponse
}
