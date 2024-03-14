package product

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type ProductService struct {
	repository Repository
}

type Service interface {
	Create(ctx context.Context, req CreateProductPayload) Response
	List(ctx context.Context, req ListProductPayload) Response
	Update(ctx context.Context, req UpdateProductPayload) Response
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

func (s *ProductService) List(ctx context.Context, req ListProductPayload) Response {
	var listProductsResponse []ProductResponse

	products, pagination, err := s.repository.List(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNoRecords
		}
		slog.Error("error fetching products list: %v", err)
		return ErrorInternal
	}

	if len(products) == 0 {
		return ErrorNoRecords
	}

	for i := range products {
		listProductsResponse = append(listProductsResponse, CreateProductResponse(products[i]))
	}

	resp := SuccessListResponse
	resp.Data = listProductsResponse
	resp.Meta = pagination

	return resp
}

func (s *ProductService) Update(ctx context.Context, req UpdateProductPayload) Response {
	// get product to authorize
	oldP, err := s.repository.GetByUUID(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNoRecords
		}
		slog.Error("error fetching product: %v", err)
		return ErrorInternal
	}

	if oldP.User.ID != req.UserID {
		return ErrorUnauthorized
	}

	newP := &Product{
		UUID:          req.ProductUID,
		Name:          req.Name,
		Price:         req.Price,
		ImageURL:      req.ImageURL,
		Condition:     req.Condition,
		Tags:          req.Tags,
		IsPurchasable: req.IsPurchasable,
		User: user.User{
			ID: req.UserID,
		},
	}

	err = s.repository.Update(ctx, newP)
	if err != nil {
		slog.Error("error patching products list: %v", err)
		return ErrorInternal
	}

	return SuccessPatchResponse
}
