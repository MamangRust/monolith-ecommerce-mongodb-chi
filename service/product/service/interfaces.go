package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-product/repository"
	salemongo "github.com/MamangRust/monolith-ecommerce-pkg/database/mongo"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, error)
	FindActive(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, error)
	FindTrashed(ctx context.Context, req *requests.FindAllProduct) ([]*repository.ProductResult, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllProductByMerchant) ([]*repository.ProductResult, *int, error)
	FindByCategory(ctx context.Context, req *requests.FindAllProductByCategory) ([]*repository.ProductResult, *int, error)
	FindByID(ctx context.Context, product_id int) (*repository.ProductResult, error)
}

type ProductCommandService interface {
	Create(ctx context.Context, req *requests.CreateProductRequest) (*salemongo.Product, error)
	Update(ctx context.Context, req *requests.UpdateProductRequest) (*salemongo.Product, error)
	UpdateProductCountStock(ctx context.Context, productID int, stock int) (*salemongo.Product, error)
	AdjustProductStock(ctx context.Context, productID int, delta int, operationID string) (*salemongo.Product, error)
	Trash(ctx context.Context, productID int) (interface{}, error)
	Restore(ctx context.Context, productID int) (interface{}, error)
	DeletePermanent(ctx context.Context, productID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}
