package product

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	bankaccount "github.com/citadel-corp/shopifyx-marketplace/internal/bank_account"
	"github.com/citadel-corp/shopifyx-marketplace/internal/user"
)

type ProductService struct {
	repository     Repository
	userRepository user.Repository
	bankRepository bankaccount.Repository
}

type Service interface {
	Create(ctx context.Context, req CreateProductPayload) Response
	List(ctx context.Context, req ListProductPayload) Response
	Update(ctx context.Context, req UpdateProductPayload) Response
	Get(ctx context.Context, req GetProductPayload) Response
	Purchase(ctx context.Context, req PurchaseProductPayload) Response
	UpdateStock(ctx context.Context, req UpdateStockPayload) Response
	Delete(ctx context.Context, req DeleteProductPayload) Response
}

func NewService(repository Repository, userRepository user.Repository, bankRepository bankaccount.Repository) Service {
	return &ProductService{
		repository:     repository,
		userRepository: userRepository,
		bankRepository: bankRepository,
	}
}

func (s *ProductService) Create(ctx context.Context, req CreateProductPayload) Response {
	serviceName := "product.Create"

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
		slog.Error(serviceName + ": " + err.Error())
		return ErrorInternal
	}

	return SuccessCreateResponse
}

func (s *ProductService) List(ctx context.Context, req ListProductPayload) Response {
	serviceName := "product.List"
	var listProductsResponse []ProductResponse

	products, pagination, err := s.repository.List(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNoRecords
		}
		slog.Error("%s: error fetching products list: %v", serviceName, err)
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
	serviceName := "product.Update"

	// get product to authorize
	oldP, err := s.repository.GetByUUID(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error fetching product: %v", serviceName, err)
		return ErrorInternal
	}

	if oldP.User.ID != req.UserID {
		return ErrorForbidden
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
		slog.Error("%s: error patching products list: %v", serviceName, err)
		return ErrorInternal
	}

	return SuccessPatchResponse
}

func (s *ProductService) Get(ctx context.Context, req GetProductPayload) Response {
	serviceName := "product.Get"

	product, err := s.repository.GetByUUID(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error fetching product: %v", serviceName, err)
		return ErrorInternal
	}

	user, err := s.userRepository.GetByID(ctx, product.User.ID)
	if err != nil {
		slog.Error("%s: error fetching product: %v", serviceName, err)
		return ErrorInternal
	}

	accts, err := s.bankRepository.List(ctx, product.User.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("%s: error fetching product: %v", serviceName, err)
			return ErrorInternal
		}
	}

	data := ProductDetailResponse{
		Product: ProductResponse{
			UUID:          product.UUID,
			Name:          product.Name,
			ImageURL:      product.ImageURL,
			Stock:         product.Stock,
			Condition:     product.Condition,
			Tags:          product.Tags,
			IsPurchasable: product.IsPurchasable,
			Price:         product.Price,
			PurchaseCount: product.PurchaseCount,
		},
		Seller: CreateSellerResponse(user, accts),
	}

	resp := SuccessGetResponse
	resp.Data = data

	return resp
}

func (s *ProductService) Purchase(ctx context.Context, req PurchaseProductPayload) Response {
	serviceName := "product.Purchase"

	// get product
	product, err := s.repository.GetByUUID(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error fetching product: %v", serviceName, err)
		return ErrorInternal
	}

	if !product.IsPurchasable {
		return ErrorNotPurchasable
	}

	if product.Stock < req.Quantity {
		return ErrorInsufficientStock
	}

	req.SellerID = product.User.ID

	// check bank account validity
	acct, err := s.bankRepository.GetByUUID(ctx, req.BankAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorBadRequest
		}
		slog.Error("%s: error fetching bank: %v", serviceName, err)
		return ErrorInternal
	}

	if acct.User.ID != req.SellerID {
		slog.Error("%s: bank does not belong to product owner: %v", serviceName, err)
		return ErrorBadRequest
	}

	err = s.repository.Purchase(ctx, req)
	if err != nil {
		slog.Error("%s: error purchasing product: %v", serviceName, err)
		return ErrorInternal
	}

	return SuccessPurchaseResponse
}

func (s *ProductService) UpdateStock(ctx context.Context, req UpdateStockPayload) Response {
	serviceName := "product.UpdateStock"

	p, err := s.repository.GetByUUID(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error fetching product: %v", serviceName, err)
		return ErrorInternal
	}

	if p.User.ID != req.UserID {
		return ErrorForbidden
	}

	product := &Product{
		UUID:  req.ProductUID,
		Stock: req.Stock,
	}
	err = s.repository.Patch(ctx, product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error updating stock: %v", serviceName, err)
		return ErrorInternal
	}

	return SuccessUpdateStockResponse
}

func (s *ProductService) Delete(ctx context.Context, req DeleteProductPayload) Response {
	serviceName := "product.Delete"

	product, err := s.repository.GetByUUID(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error fetching product: %v", serviceName, err)
		return ErrorInternal
	}

	if product.User.ID != req.UserID {
		return ErrorForbidden
	}

	err = s.repository.Delete(ctx, req.ProductUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorNotFound
		}
		slog.Error("%s: error deleting product: %v", serviceName, err)
		return ErrorInternal
	}

	return SuccessDeleteResponse
}
