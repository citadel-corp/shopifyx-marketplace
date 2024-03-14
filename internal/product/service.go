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

	return SuccessCreateResponse
}

func (s *ProductService) List(ctx context.Context, req ListProductPayload) ([]ProductResponse, *response.Pagination, Response) {
	var listProductsResponse []ProductResponse

	products, err := s.repository.List(ctx, req)
	if err != nil {
		if len(products) == 0 {
			return listProductsResponse, nil, ErrorNoRecords
		}
	}

	for i := range products {
		listProductsResponse = append(listProductsResponse, CreateProductResponse(products[i]))
	}

	pagination := &response.Pagination{
		Limit:  req.Limit,
		Offset: req.Offset,
		Total:  len(listProductsResponse),
	}

	return listProductsResponse, pagination, SuccessListResponse
}
