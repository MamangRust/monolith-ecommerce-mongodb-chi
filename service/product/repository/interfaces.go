package repository

import (
	"context"

	categoryadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/merchant"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// CategoryQueryRepository delegates to the shared category adapter.
type CategoryQueryRepository = categoryadapter.QueryRepository

// MerchantQueryRepository delegates to the shared merchant adapter.
type MerchantQueryRepository = merchantadapter.QueryRepository

type ProductResult struct {
	ProductID    int32
	MerchantID   int32
	CategoryID   int32
	Name         string
	Description  *string
	Price        int32
	CountInStock int32
	Brand        *string
	Weight       *int32
	Rating       *float64
	SlugProduct  *string
	ImageProduct *string
	CreatedAt    *string
	UpdatedAt    *string
	DeletedAt    *string
	TotalCount   int64
}

type ProductQueryRepository interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error)
	FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error)
	FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*ProductResult, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*ProductResult, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*ProductResult, error)
	FindByID(ctx context.Context, product_id int) (*ProductResult, error)
	FindByIDTrashed(ctx context.Context, product_id int) (*ProductResult, error)
}

type ProductCommandRepository interface {
	Create(ctx context.Context, request *requests.CreateProductRequest) (*salemongo.Product, error)
	Update(ctx context.Context, request *requests.UpdateProductRequest) (*salemongo.Product, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*salemongo.Product, error)
	AdjustProductStock(ctx context.Context, product_id int, delta int, operationID string) (*salemongo.Product, error)
	Trash(ctx context.Context, product_id int) (*salemongo.Product, error)
	Restore(ctx context.Context, product_id int) (*salemongo.Product, error)
	DeletePermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
