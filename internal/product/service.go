package product

import (
	"context"
	"log/slog"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type ProductService struct {
	repository Repository
}

type Service interface {
	Create(ctx context.Context, req CreateProductPayload) Response
	List(ctx context.Context, req ListProductPayload) ([]ProductResponse, *response.Pagination, Response)
}

func NewService(repository Repository) Service {
	return &ProductService{repository: repository}
}

func (s *ProductService) Create(ctx context.Context, req CreateProductPayload) Response {
	product := &Product{
		Name:          req.Name,
		ImageURL:      req.ImageURL,
		Stock:         req.Stock,
		Condition:     req.Condition,
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

	return SuccessCreateResponse
}

func (s *ProductService) List(ctx context.Context, req ListProductPayload) ([]ProductResponse, *response.Pagination, Response) {
	var listProductsResponse []ProductResponse

	products, pagination, err := s.repository.List(ctx, req)
	if err != nil {
		slog.Error("error fetching products list: %v", err)
		return nil, nil, ErrorInternal
	}

	if len(products) == 0 {
		return listProductsResponse, nil, ErrorNoRecords
	}

	for i := range products {
		listProductsResponse = append(listProductsResponse, CreateProductResponse(products[i]))
	}

	return listProductsResponse, pagination, SuccessListResponse
}
