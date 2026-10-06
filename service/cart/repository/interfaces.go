package repository

import (
	"context"

	productadapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/product"
	useradapter "github.com/MamangRust/monolith-ecommerce-pkg/adapter/user"

	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
)

// CartResult is the result type for paginated cart queries.
type CartResult struct {
	CartID     int32
	UserID     int32
	ProductID  int32
	Name       string
	Price      int32
	Image      string
	Quantity   int32
	Weight     int32
	CreatedAt  string
	UpdatedAt  string
	TotalCount int64
}

// CartCreateResult is the result type for cart creation.
type CartCreateResult struct {
	CartID    int32
	UserID    int32
	ProductID int32
	Name      string
	Price     int32
	Image     string
	Quantity  int32
	Weight    int32
	CreatedAt string
	UpdatedAt string
}

type CartQueryRepository interface {
	FindCarts(ctx context.Context, req *requests.FindAllCarts) ([]*CartResult, error)
}

type CartCommandRepository interface {
	CreateCart(ctx context.Context, req *requests.CartCreateRecord) (*CartCreateResult, error)
	DeletePermanent(ctx context.Context, req *requests.DeleteCartRequest) (bool, error)
	DeleteAllPermanently(ctx context.Context, req *requests.DeleteAllCartRequest) (bool, error)
}

// ProductQueryRepository delegates to the shared product adapter.
type ProductQueryRepository = productadapter.QueryRepository

// UserQueryRepository delegates to the shared user adapter.
type UserQueryRepository = *useradapter.Repository

type ProductResult struct {
	ProductID    int32
	Name         string
	Price        int32
	CountInStock int32
	ImageProduct *string
	Weight       *int32
}

type UserResult struct {
	UserID int32
}
